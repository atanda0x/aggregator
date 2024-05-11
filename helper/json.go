package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json res: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
