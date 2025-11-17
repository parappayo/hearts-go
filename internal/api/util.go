package api

import (
	"encoding/json"
	"net/http"
	"log"
)

func WriteResponse[T any](w http.ResponseWriter, data T) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: failed to marshal response", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(body))
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
