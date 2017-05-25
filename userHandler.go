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

type userException struct {
	Message string
	Err     error
}

const userCollection = "users"

func valuFromContext(r *http.Request) (*config.Config, *gorm.DB) {
	conf := r.Context().Value(configuration).(*config.Config)
	DB := r.Context().Value(psqlDB).(*gorm.DB)

	return conf, DB
}

func decodeBody(requestBody io.Reader) (user, userException) {
	decoder := json.NewDecoder(requestBody)
	var (
		u  user
		ue userException
	)

	err := decoder.Decode(&u)
	if err != nil {
		ue.Err = err
		ue.Message = "Error on body decode"
		return u, ue
	}
	return u, ue
}

func getUser(w http.ResponseWriter, r *http.Request) {
	_, db := valuFromContext(r)
	userID := strings.Title(chi.URLParam(r, "id"))

	u := &user{}
	err := db.First(&u, userID).Error
	if err != nil {
		ue := userException{
			Message: "Error on retrieving user",
			Err:     err,
		}
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	u, ue := decodeBody(r.Body)
	if ue.Err != nil {
		helper.WriteJSON(w, ue, http.StatusBadRequest)
	}

	_, db := valuFromContext(r)
	if u.userAlreadyExists(db) {
		ue.Message = "User name or email already exists"
		ue.Err = nil
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}

	err := db.Create(&u).Error
	if err != nil {
		ue.Message = "Error on DB Create "
		ue.Err = err
		helper.WriteJSON(w, ue, http.StatusInternalServerError)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	u, ue := decodeBody(r.Body)
	if ue.Err != nil {
		helper.WriteJSON(w, ue, http.StatusBadRequest)
	}

	_, db := valuFromContext(r)

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
