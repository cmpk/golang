package entry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/model"
	"backend/singleton"
)

const RECORD_URL = "/api/record"

var RecordHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	sessionManager := singleton.GetSessionManager()
	session := sessionManager.SessionStart(w, req)
	uid, _ := session.Get("uid")

	connection := model.CreateConnection()
	defer connection.Close()

	var err error

	records, err := model.GetRecordsByUid(connection, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(ret))
})
