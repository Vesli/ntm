package helper

import (
	"encoding/json"
	"net/http"
)

// WriteJSON is use to write a JSON response (error and success) from object on request
func WriteJSON(w http.ResponseWriter, obj interface{}, responseStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("error on encoding data"))
	}
}
