package singleton

import (
	"backend/session"
	"os"
)

var sessionManager *session.Manager

func GetSessionManager() *session.Manager {
	if sessionManager == nil {
		var err error
		sessionManager, err = session.CreateManager(os.Getenv("APP_NAME"), 3600)
		if err != nil {
			panic(err)
		}
	}
	return sessionManager
}
