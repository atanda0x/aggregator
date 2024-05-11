package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// Respond the requet from client in JSON
func resWithJON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON res: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "appliction/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
