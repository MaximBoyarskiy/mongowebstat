package mongowebstat

import (
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
)

type Node struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Http bool   `json:"http"`
}

var (
	nodes []Node = make([]Node, 0)
)

func (n *Node) connString() string {
	if n.Http {
		return n.Host + "/_status"
	}
	return n.Host

}

func (n *Node) hostName() string {
	if n.Name == "" {
		return n.Host
	}
	return n.Name
}

func (n *Node) GetMongoStatRaw(c chan NodeResponse) {
	defer func() {
		if r := recover(); r != nil {
			//TODO
			//Fix closed channel issue
			//l.Print("If it is about closed channel nevermind: ", r)
			if fmt.Sprintf("%v", r) != "runtime error: send on closed channel" {
				panic(r)
			}
		}
	}()

	var result NodeResponse
	nn := n.hostName()
	nc := n.connString()
	debug.Print(nn, " node.GetMongoStatRaw")
	if n.Http {
		debug.Print(nn, " n.Http")
		if r, err := http.Get(nc); err == nil {
			if b, err := ioutil.ReadAll(r.Body); err == nil {
				stats, err := byteArrayToMongoStatRaw(b)
				if err != nil {
					debug.Print(nn, " byteArrayToMongoStatRaw ", err.Error())
				}
				result = NodeResponse{NodeStatus{err: err, stats: stats}, nn}
			} else {
				debug.Print(nn, " ioutil.ReadAll err: ", err.Error())
				result = NodeResponse{NodeStatus{err: err, stats: &MongoStatRaw{}}, nn}
			}
			r.Body.Close()
		} else {
			result = NodeResponse{NodeStatus{err: err, stats: &MongoStatRaw{}}, nn}
		}
	} else {
		debug.Print(nn, "!n.Http")
		//if session, err := mgo.Dial(nc); err == nil {
		if session, err := mgo.DialWithTimeout(nc, timeOutMgo); err == nil {
			defer session.Close()
			data := map[string]interface{}{}
			if err = session.Run("serverStatus", data); err == nil {
				stats, err := mapToMongoStatRaw(map[string]interface{}{"serverStatus": data})
				if err != nil {
					debug.Print(nn, " toMongoStatRaw retured err: ", err.Error())
				}
				result = NodeResponse{NodeStatus{err: err, stats: stats}, nn}
			} else {
				debug.Print(nn, " session.Run err: ", err.Error())
				result = NodeResponse{NodeStatus{err: err, stats: &MongoStatRaw{}}, nn}
			}
		} else {
			debug.Print(nn, " mgo.Dial(nc) err: ", err.Error())
			result = NodeResponse{NodeStatus{err: err, stats: &MongoStatRaw{}}, nn}
		}
	}
	c <- result
}
