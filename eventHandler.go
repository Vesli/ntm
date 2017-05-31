package main

import (
	"net/http"

	"github.com/vesli/ntm/helper"
)

/*
	TODO CreateEvent: UserID passed to the foreignKey has to be from the connected user
	and not from JSON.
*/
func createEvent(w http.ResponseWriter, r *http.Request) {
	var e *Event

	err := helper.DecodeBody(e, r.Body)
	if err != nil {
		helper.WriteJSON(w, err, http.StatusBadRequest)
	}

	_, db := valueFromContext(r)
	err = db.Create(&e).Error
	if err != nil {
		helper.WriteJSON(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, e, http.StatusOK)
}
