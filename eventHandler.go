package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vesli/ntm/helper"
)

func decodeEventBody(requestBody io.Reader) (Event, helper.ResponseException) {
	decoder := json.NewDecoder(requestBody)
	var (
		e  Event
		re helper.ResponseException
	)

	err := decoder.Decode(&e)
	if err != nil {
		re.Err = err
		re.Message = "Error on body decode"
		return e, re
	}
	return e, re
}

/*
	TODO CreateEvent: UserID passed to the foreignKey has to be from the connected user
	and not from JSON.
*/
func createEvent(w http.ResponseWriter, r *http.Request) {
	e, re := decodeEventBody(r.Body)
	if re.Err != nil {
		helper.WriteJSON(w, re, http.StatusBadRequest)
	}

	_, db := valueFromContext(r)
	err := db.Create(&e).Error
	if err != nil {
		re.Message = "Error on DB Create"
		re.Err = err
		helper.WriteJSON(w, re, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, e, http.StatusOK)
}
