package mongowebstat

import (
	"encoding/json"
	"net/http"
)

func toJsonResponse(w http.ResponseWriter, r interface{}) {
	b, err := json.Marshal(r)
	if err != nil {
		l.Println("error: ", err.Error())
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.Write(b)
}
