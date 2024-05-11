package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx err: %v\n", msg)
	}

	type errResonse struct {
		Error string `json:"error"`
	}

	ResWithJSON(w, code, errResonse{
		Error: msg,
	})
}

func ResWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json res: %v\n", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
