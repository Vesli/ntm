package helper

import (
	"encoding/json"
	"net/http"
)

// WriteJSON is use to write a JSON response from object on request
func WriteJSON(w http.ResponseWriter, obj interface{}, responseStatus int) {
	jData, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(responseStatus)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(jData)
}
