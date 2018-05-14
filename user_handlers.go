package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ntm/helper"
)

type token struct {
	AccessToken string `json:"access_token"`
}

func (s *Service) checkFBToken(token string) error {
	client := http.Client{}

	urlParams := url.Values{}
	urlParams.Add("access_token", token)
	urlParams.Add("fields", s.Conf.FBFields)

	req, err := http.NewRequest("POST", s.Conf.FBURL, strings.NewReader(urlParams.Encode()))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//TO DO RESP
	return err
}

// RegisterOrLoggin manage login via FB, refresh token etc.
func (s *Service) RegisterOrLoggin(w http.ResponseWriter, r *http.Request) {
	t := &token{}

	log.Println("t: ", t)
	err := helper.DecodeBody(&t, r.Body)
	if err != nil {
		// TODO does the writeJSON correctly write errors back to client?
		helper.WriteJSON(w, err, http.StatusBadRequest)
		return
	}

	/*
	  Request fields as image /link, Name, age.
	*/

	s.checkFBToken(t.AccessToken)
}
