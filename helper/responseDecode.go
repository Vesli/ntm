package helper

import (
	"encoding/json"
	"io"
)

// DecodeBody that help to directly fill every obj from JSON
func DecodeBody(v interface{}, requestBody io.Reader) error {
	decoder := json.NewDecoder(requestBody)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}
