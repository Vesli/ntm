package main

import (
	"net/http"

	"github.com/ntm/helper"
)

type token struct {
	AccessToken string `json:"access_token"`
}

// RegisterOrLoggin manage login via FB, refresh token etc.
func (s *Service) RegisterOrLoggin(w http.ResponseWriter, r *http.Request) {
	t := &token{}

	err := helper.DecodeBody(&t, r.Body)
	if err != nil {
		// TODO does the writeJSON correctly write errors?
		helper.WriteJSON(w, err, http.StatusBadRequest)
	}

}
