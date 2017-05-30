package main

/*
	The userHandler file is used as a user manager.
	Here you have the main user actions.
	Need to rewrite the requests response.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pressly/chi"
	"github.com/vesli/ntm/config"
	"github.com/vesli/ntm/helper"
)

// Token structure for passing access token from response to the facebook graph API
type Token struct {
	AccessToken string `json:"access_token"`
}

func valueFromContext(r *http.Request) (*config.Config, *gorm.DB) {
	conf := r.Context().Value(configuration).(*config.Config)
	DB := r.Context().Value(psqlDB).(*gorm.DB)

	return conf, DB
}

func decodeBody(requestBody io.Reader) (*Token, helper.ResponseException) {
	var (
		t  Token
		re helper.ResponseException
	)

	decoder := json.NewDecoder(requestBody)
	err := decoder.Decode(&t)
	if err != nil {
		re.Err = err
		re.Message = "Error on decode"
		return nil, re
	}

	return &t, re
}

/* Facebook login from URL */
func getUserFromToken(t *Token, conf *config.Config, re helper.ResponseException) (*User, helper.ResponseException) {
	var u *User

	urlParams := make(url.Values)
	urlParams.Add("access_token", t.AccessToken)
	urlParams.Add("fields", conf.FBParams)

	ret, err := http.Get(fmt.Sprintf("%s?%s", conf.FBURL, urlParams.Encode()))
	if err != nil {
		re.Err = err
		re.Message = "Error on get params"
		return nil, re
	}

	decoder := json.NewDecoder(ret.Body)
	err = decoder.Decode(&u)
	if err != nil {
		re.Err = err
		re.Message = "Error on decode"
		return nil, re
	}
	return u, re
}

func getUserFromDB(w http.ResponseWriter, r *http.Request) {
	_, db := valueFromContext(r)
	userID := strings.Title(chi.URLParam(r, "id"))

	u := &User{}
	err := db.First(&u, userID).Error
	if err != nil {
		re := helper.ResponseException{
			Message: "Error on retrieving user",
			Err:     err,
		}
		helper.WriteJSON(w, re, http.StatusBadRequest)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func registerAndLogginUser(w http.ResponseWriter, r *http.Request) {
	conf, db := valueFromContext(r)

	t, re := decodeBody(r.Body)
	if re.Err != nil {
		helper.WriteJSON(w, re, http.StatusBadRequest)
		return
	}

	u, re := getUserFromToken(t, conf, re)
	if re.Err != nil {
		helper.WriteJSON(w, re, http.StatusBadRequest)
		return
	}
	if !u.userAlreadyExists(db) {
		err := db.Create(&u).Error
		if err != nil {
			re.Message = "Error on DB Create "
			re.Err = err
			helper.WriteJSON(w, re, http.StatusInternalServerError)
			return
		}
	}

	helper.WriteJSON(w, u, http.StatusOK)
}
