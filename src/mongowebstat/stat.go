package mongowebstat

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	stats              Stats
	period             = 1000 * time.Millisecond //It should be the same as in angularjs app setting.
	timeOutMgo         = 800 * time.Millisecond  //Should be a bit less then period otherwise we didn't get appropriate mgo error. However might be equal period
	unavailable        = "No data available"
	na                 = "N/A"
	doesntresponse     = fmt.Errorf("Resourse doesn't response for %v;", period)
	MongoStatRawFields []MongoStatRawField
	MongoStatFields    []MongoStatField
	True               bool = true
	False              bool = false
)

const (
	MaxInt32 = int64(^uint32(0) >> 1)
)

type MongoStat struct {
	Insert  int32 `op:"diff" args:"Insert"`
	Query   int32 `op:"diff" args:"Query"`
	Update  int32 `op:"diff" args:"Update"`
	Delete  int32 `op:"diff" args:"Delete"`
	Getmore int32 `op:"diff" args:"Getmore"`
	Command int32 `op:"diff" args:"Command"`

	ReplInsert  int32 `op:"diff" args:"ReplInsert"`
	ReplQuery   int32 `op:"diff" args:"ReplQuery"`
	ReplUpdate  int32 `op:"diff" args:"ReplUpdate"`
	ReplDelete  int32 `op:"diff" args:"ReplDelete"`
	ReplGetmore int32 `op:"diff" args:"ReplGetmore"`
	ReplCommand int32 `op:"diff" args:"ReplCommand"`

	Flushes int32 `op:"diff" args:"Flushes"`

	Mapped    int32 `op:"now" args:"Mapped" cond:"!isMongos"`
	Virtual   int32 `op:"now" args:"Virtual"`
	Resident  int32 `op:"now" args:"Resident"`
	NonMapped int32 `op:"diff" args:"Virtual,Mapped" cond:"!isMongos"`
	Faults    int32 `op:"diff" args:"Faults"`
	//TODO
	//Locks   int32 `op:"locks" args:"-"`
	IdxMiss float64 `op:"percent" args:"Accesses,Misses" cond:"!isMongos"`

	QueueReaders  int32 `op:"now" args:"QueueReaders"`
	QueueWriters  int32 `op:"now" args:"QueueWriters"`
	ActiveReaders int32 `op:"now" args:"ActiveReaders"`
	ActiveWriters int32 `op:"now" args:"ActiveWriters"`

	NetIn  int32 `op:"diff" args:"BytesIn"`
	NetOut int32 `op:"diff" args:"BytesOut"`

	Connections int32 `op:"now" args:"Connections"`
	//TODO
	//Repl string `op:"repl" args:"-"`
	//Time string `op:"time" args:"Time"`

	Error string
}

func NewMongoStat() *MongoStat {
	return &MongoStat{
		Insert:  -1,
		Query:   -1,
		Update:  -1,
		Delete:  -1,
		Getmore: -1,
		Command: -1,

		ReplInsert:  -1,
		ReplQuery:   -1,
		ReplUpdate:  -1,
		ReplDelete:  -1,
		ReplGetmore: -1,
		ReplCommand: -1,

		Flushes: -1,

		Mapped:    -1,
		Virtual:   -1,
		Resident:  -1,
		NonMapped: -1,
		Faults:    -1,

		//Locks:   -1,
		IdxMiss: -1,

		QueueReaders:  -1,
		QueueWriters:  -1,
		ActiveReaders: -1,
		ActiveWriters: -1,

		NetIn:  -1,
		NetOut: -1,

		Connections: -1,

		//Repl: "",
		//Time: "",

		Error: "",
	}
}

func (m *MongoStat) String() string {
	return fmt.Sprintf("MongoStat: %v, %v, Error: %v\n", m.Insert, m.Query, m.Error)
}

type MongoStatRaw struct {
	IsMongos        *bool
	ShardCursorType map[string]interface{} `jsonq_path:"shardCursorType"`
	Process         *string                `jsonq_path:"process"`

	Insert  *int32 `jsonq_path:"opcounters,insert"`
	Query   *int32 `jsonq_path:"opcounters,query"`
	Update  *int32 `jsonq_path:"opcounters,update"`
	Delete  *int32 `jsonq_path:"opcounters,delete"`
	Getmore *int32 `jsonq_path:"opcounters,getmore"`
	Command *int32 `jsonq_path:"opcounters,command"`

	ReplInsert  *int32 `jsonq_path:"opcountersRepl,insert"`
	ReplQuery   *int32 `jsonq_path:"opcountersRepl,query"`
	ReplUpdate  *int32 `jsonq_path:"opcountersRepl,update"`
	ReplDelete  *int32 `jsonq_path:"opcountersRepl,delete"`
	ReplGetmore *int32 `jsonq_path:"opcountersRepl,getmore"`
	ReplCommand *int32 `jsonq_path:"opcountersRepl,command"`

	Flushes *int32 `jsonq_path:"backgroundFlushing,flushes"`

	MemSupported *bool  `jsonq_path:"mem,supported"`
	Mapped       *int32 `jsonq_path:"mem,mapped"`
	Virtual      *int32 `jsonq_path:"mem,virtual"`
	Resident     *int32 `jsonq_path:"mem,resident"`

	Faults *int32 `jsonq_path:"extra_info,page_faults"`

	Locks map[string]interface{} `jsonq_path:"locks"`

	TotalTime *int64 `jsonq_path:"globalLock,totalTime"`
	LockTime  *int64 `jsonq_path:"globalLock,lockTime"`

	Accesses *int32 `jsonq_path:"indexCounters,accesses"` //Mongodb 2.0.x `jsonq_path:"indexCounters,btree,accesses"`
	Misses   *int32 `jsonq_path:"indexCounters,misses"`   //Mongodb 2.0.x `jsonq_path:"indexCounters,btree,misses"`

	QueueReaders *int32 `jsonq_path:"globalLock,currentQueue,readers"`
	QueueWriters *int32 `jsonq_path:"globalLock,currentQueue,writers"`

	ActiveReaders *int32 `jsonq_path:"globalLock,activeClients,readers"`
	ActiveWriters *int32 `jsonq_path:"globalLock,activeClients,writers"`

	BytesIn  *int32 `jsonq_path:"network,bytesIn"`
	BytesOut *int32 `jsonq_path:"network,bytesOut"`

	Connections *int32 `jsonq_path:"connections,current"`

	//TODO
	//Repl map[string]interface{} `jsonq_path:"repl"`
	//Time *int64 `jsonq_path:"localTime"`
}

type mapNodeStatus map[string]*NodeStatus

func (s *mapNodeStatus) init() {
	*s = make(mapNodeStatus)
}

type NodeStatus struct {
	stats *MongoStatRaw
	err   error
}

type Stats struct {
	sync.Mutex
	now  mapNodeStatus
	prev mapNodeStatus
	data map[string]*MongoStat
}

func (s *Stats) updateNodeData(node string) {
	debug.Print("-> ", node, " updateNodeData")
	now := s.now[node].stats
	if now.ShardCursorType != nil || (now.Process != nil && *now.Process == "mongos") {
		now.IsMongos = &True
	} else {
		now.IsMongos = &False
	}
	for _, field := range MongoStatFields {
		nv := reflect.ValueOf(s.now[node].stats).Elem()
		pv := reflect.ValueOf(s.prev[node].stats).Elem()
		dv := reflect.ValueOf(s.data[node]).Elem()
		if (field.cond == "!isMongos" && *now.IsMongos) ||
			(field.cond == "isMongos" && !*now.IsMongos) {
			continue
		}
		switch field.op {
		case "percent":
			var xp, xn, yp, yn reflect.Value
			difference := (float64)(-1)
			xn, xp = nv.FieldByName(field.args[1]), pv.FieldByName(field.args[1])
			yn, yp = nv.FieldByName(field.args[0]), pv.FieldByName(field.args[0])
			d := dv.FieldByName(field.name)
			if !(xn.IsNil() || xp.IsNil()) && !(yn.IsNil() || yp.IsNil()) {
				x := (float64)(xn.Elem().Int() - xp.Elem().Int())
				y := (float64)(yn.Elem().Int() - yp.Elem().Int())
				if y == 0 {
					difference = (float64)(0)
				} else {
					difference = x / y
					difference = (float64)((int64)(difference*1000)) / 10
				}
			}
			d.SetFloat(difference)
		case "diff":
			debug.Print("diff ", field.name)
			var p, n reflect.Value
			difference := (int64)(-1)
			n = nv.FieldByName(field.args[0])
			if len(field.args) == 2 {
				p = nv.FieldByName(field.args[1])
			} else {
				p = pv.FieldByName(field.args[0])
			}
			d := dv.FieldByName(field.name)
			if !(n.IsNil() || p.IsNil()) {
				minuend, subtrahend := n.Elem().Int(), p.Elem().Int()
				debug.Print("minuend, subtrahend", minuend, subtrahend)
				if minuend >= subtrahend {
					difference = (int64)(minuend - subtrahend)
				} else {
					difference = (int64)(MaxInt32 - subtrahend + minuend)
				}
			}
			d.SetInt(difference)
		case "now":
			arg := field.args[0]
			n := nv.FieldByName(arg)
			d := dv.FieldByName(field.name)
			var v int64 = -1
			if !n.IsNil() {
				v = n.Elem().Int()
			}
			debug.Printf("now %v %v", field.name, v)
			d.SetInt(v)

		}
	}

}

func (s *Stats) updateData() {
	debug.Print("-> updateData")
	stats.Lock()
	defer stats.Unlock()
	s.data = make(map[string]*MongoStat)
	for _, n := range nodes {
		node := n.hostName()
		s.data[node] = NewMongoStat()
		now, nowExist := s.now[node]
		if nowExist && (now.err != nil) {
			fmt.Printf("\nnow.err: %v;\n", now.err.Error())
			//print(s.data[node].String())
			s.data[node].Error = now.err.Error()
			continue
		}
		prev, prevExist := s.prev[node]
		if nowExist && prevExist && prev.err == nil {
			s.updateNodeData(node)
		} else {
			s.data[node].Error = unavailable
		}
	}
}

type NodeResponse struct {
	NodeStatus
	host string
}

type MongoStatRawField struct {
	name      string
	kind      reflect.Kind
	jsonqPath []string
}

type MongoStatField struct {
	name string
	op   string
	args []string
	cond string
}

func init() {
	stats = Stats{now: make(map[string]*NodeStatus), prev: make(map[string]*NodeStatus), data: make(map[string]*MongoStat)}

	MongoStatRawType := reflect.TypeOf(MongoStatRaw{})
	l := MongoStatRawType.NumField()
	MongoStatRawFields = make([]MongoStatRawField, l)
	for i := 0; i < l; i++ {
		fv := MongoStatRawType.Field(i)
		jsonq_path := strings.Split(fv.Tag.Get("jsonq_path"), ",")
		if jsonq_path[0] != "" {
			jsonq_path = append([]string{"serverStatus"}, jsonq_path...)
		}
		var fk reflect.Kind
		if fv.Type.Kind() == reflect.Ptr {
			fk = fv.Type.Elem().Kind()
		} else {
			fk = fv.Type.Kind()
		}
		MongoStatRawFields[i] = MongoStatRawField{fv.Name, fk, jsonq_path}
	}

	MongoStatType := reflect.TypeOf(MongoStat{})
	l = MongoStatType.NumField()
	MongoStatFields = make([]MongoStatField, l)
	for i := 0; i < l; i++ {
		fv := MongoStatType.Field(i)
		op := fv.Tag.Get("op")
		args := strings.Split(fv.Tag.Get("args"), ",")
		cond := fv.Tag.Get("cond")
		MongoStatFields[i] = MongoStatField{fv.Name, op, args, cond}
	}
}

func puller() {
	t := time.Tick(period)
	var result mapNodeStatus
	for {
		debug.Print("New cycle!\n\n")
		c := make(chan NodeResponse)
		for _, n := range nodes {
			go func(n Node, c chan NodeResponse) {
				n.GetMongoStatRaw(c)
			}(n, c)
		}
		result.init()
		debug.Print("Waiting for results...")
	inside:
		for {
			select {
			case resp := <-c:
				result[resp.host] = &NodeStatus{resp.stats, resp.err}
				debug.Printf("resp.host %v:\nresp.stats: %v\nresp.stats.Query: %v", resp.host, resp.stats, resp.stats.Query)
			case _ = <-t:
				debug.Print("Timeout!")
				close(c)
				for _, n := range nodes {
					node := n.hostName()
					debug.Print("Time to stat result for ", node)

					if n, ok := stats.now[node]; ok {
						debug.Print("now is here")
						stats.prev[node] = n
					} else {
						debug.Print("now is not here")
						stats.prev[node] = &NodeStatus{err: doesntresponse, stats: &MongoStatRaw{}}
					}
					if r, ok := result[node]; ok {
						debug.Print("result is here")
						stats.now[node] = r
					} else {
						debug.Print("result is not here")
						stats.now[node] = &NodeStatus{err: doesntresponse, stats: &MongoStatRaw{}}
					}
				}
				stats.updateData()
				break inside
			}
		}
	}
}
