package mongowebstat

import (
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func ServerStart(httpPtr *string) {

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/config", config)
	http.HandleFunc("/stats", getStats)
	http.HandleFunc("/nodes", getNodes)

	l.Print("Server is started.")
	go puller()
	if err := http.ListenAndServe(*httpPtr, nil); err != nil {
		l.Fatal("Server error: ", err.Error())
	}
}
