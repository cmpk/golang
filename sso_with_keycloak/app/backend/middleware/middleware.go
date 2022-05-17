package middleware

import (
	"net/http"
)

func SetAccessControlOnHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// log.Print("===== START : SetAccessControlOnHeaderMiddleware")
		handler.ServeHTTP(w, req)

		// w.Header().Set("Access-Control-Allow-Origin", os.Getenv("APP_URL"))
		// // w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,HEAD,OPTIONS")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// log.Print("===== END : SetAccessControlOnHeaderMiddleware")
	})
}
