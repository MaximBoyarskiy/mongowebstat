package mongowebstat

import (
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Start() {

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./mongowebstat/static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/config", config)
	http.HandleFunc("/stats", getStats)
	http.HandleFunc("/nodes", getNodes)

	l.Print("Server is started.")
	go puller()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		l.Fatal("Server error: ", err.Error())
	}
}
