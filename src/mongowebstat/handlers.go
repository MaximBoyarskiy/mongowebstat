package mongowebstat

import (
	"io/ioutil"
	"net/http"
)

var tmpl []byte

func init() {
	var err error
	tmpl, err = ioutil.ReadFile("./mongowebstat/templates/base.html")
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write(tmpl)
}

func config(w http.ResponseWriter, r *http.Request) {
	if err := LoadConfig(); err != nil {
		panic(err)
	}
	index(w, r)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	stats.Lock()
	toJsonResponse(w, stats.data)
	stats.Unlock()

}

func getNodes(w http.ResponseWriter, r *http.Request) {
	n := make([]struct{ Name string }, len(nodes))
	for i, v := range nodes {
		n[i] = struct{ Name string }{Name: v.hostName()}
	}
	toJsonResponse(w, n)
}
