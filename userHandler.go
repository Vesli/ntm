package main

/*
	The userHandler file is used as a user manager.
	Here you have the main user actions.
	Need to rewrite the requests response.
*/

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pressly/chi"
	"github.com/vesli/ntm/config"
	"github.com/vesli/ntm/helper"
)

func valueFromContext(r *http.Request) (*config.Config, *gorm.DB) {
	conf := r.Context().Value(configuration).(*config.Config)
	DB := r.Context().Value(psqlDB).(*gorm.DB)

	return conf, DB
}

func decodeUserBody(requestBody io.Reader) (User, helper.ResponseException) {
	decoder := json.NewDecoder(requestBody)
	var (
		u  User
		re helper.ResponseException
	)

	err := decoder.Decode(&u)
	if err != nil {
		re.Err = err
		re.Message = "Error on body decode"
		return u, re
	}
	return u, re
}

func getUser(w http.ResponseWriter, r *http.Request) {
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

func registerUser(w http.ResponseWriter, r *http.Request) {
	u, re := decodeUserBody(r.Body)
	if re.Err != nil {
		helper.WriteJSON(w, re, http.StatusBadRequest)
	}

	_, db := valueFromContext(r)
	if u.userAlreadyExists(db) {
		re.Message = "User name or email already exists"
		re.Err = nil
		helper.WriteJSON(w, re, http.StatusBadRequest)
		return
	}

	err := db.Create(&u).Error
	if err != nil {
		re.Message = "Error on DB Create "
		re.Err = err
		helper.WriteJSON(w, re, http.StatusInternalServerError)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	u, ue := decodeUserBody(r.Body)
	if ue.Err != nil {
		helper.WriteJSON(w, ue, http.StatusBadRequest)
	}

	_, db := valueFromContext(r)

	accessToken, err := u.checkCredentials(db)
	if err != nil {
		ue.Err = err
		ue.Message = "wrong password or login"
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}

	u.AccessToken = accessToken
	db.Model(&u).Update(&u)
	helper.WriteJSON(w, u, http.StatusOK)
}
