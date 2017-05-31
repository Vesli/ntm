package main

/*
	The userHandler file is used as a user manager.
	Here you have the main user actions.
	Need to rewrite the requests response.
*/

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pressly/chi"
	"github.com/vesli/ntm/config"
	"github.com/vesli/ntm/helper"
)

// Token structure, passing access token from body response to the facebook graph API
type Token struct {
	AccessToken string `json:"access_token"`
}

func valueFromContext(r *http.Request) (*config.Config, *gorm.DB) {
	conf := r.Context().Value(configuration).(*config.Config)
	DB := r.Context().Value(psqlDB).(*gorm.DB)

	return conf, DB
}

/* Facebook login from URL */
func getUserFromToken(t *Token, conf *config.Config) (*User, error) {
	u := &User{}

	urlParams := make(url.Values)
	urlParams.Add("access_token", t.AccessToken)
	urlParams.Add("fields", conf.FBParams)

	ret, err := http.Get(fmt.Sprintf("%s?%s", conf.FBURL, urlParams.Encode()))
	if err != nil {
		return nil, err
	}

	err = helper.DecodeBody(&u, ret.Body)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func getUserFromDB(w http.ResponseWriter, r *http.Request) {
	_, db := valueFromContext(r)
	userID := strings.Title(chi.URLParam(r, "id"))

	u := &User{}
	err := db.First(&u, userID).Error
	if err != nil {
		helper.WriteJSON(w, err, http.StatusBadRequest)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func registerAndLogginUser(w http.ResponseWriter, r *http.Request) {
	t := &Token{}

	re := helper.DecodeBody(&t, r.Body)
	if re != nil {
		helper.WriteJSON(w, re, http.StatusBadRequest)
		return
	}

	conf, db := valueFromContext(r)

	u, err := getUserFromToken(t, conf)
	if err != nil {
		helper.WriteJSON(w, err, http.StatusBadRequest)
		return
	}
	if !u.userAlreadyExists(db) {
		err = db.Create(&u).Error
		if err != nil {
			helper.WriteJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	helper.WriteJSON(w, u, http.StatusOK)
}
