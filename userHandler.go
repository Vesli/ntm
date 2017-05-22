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

	"github.com/vesli/ntm/config"
	"github.com/vesli/ntm/helper"
	mgo "gopkg.in/mgo.v2"
)

type userException struct {
	Message string
	Err     error
}

const userCollection = "users"

func contextFromMiddleware(r *http.Request) (*config.Config, *mgo.Session) {
	conf := r.Context().Value(configuration).(*config.Config)
	sessionC := r.Context().Value(mgoSession).(*mgo.Session)

	return conf, sessionC
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

func registerUser(w http.ResponseWriter, r *http.Request) {
	u, ue := decodeBody(r.Body)
	if ue.Err != nil {
		helper.WriteJSON(w, ue, http.StatusBadRequest)
	}

	conf, sessionC := contextFromMiddleware(r)
	c := sessionC.DB(conf.DBName).C(userCollection)

	if u.userAlreadyExists(c) {
		ue.Message = "User name or email already exists"
		ue.Err = nil
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}

	u.insertUserToDB(w, c)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	u, ue := decodeBody(r.Body)
	if ue.Err != nil {
		helper.WriteJSON(w, ue, http.StatusBadRequest)
	}

	conf, sessionC := contextFromMiddleware(r)
	c := sessionC.DB(conf.DBName).C(userCollection)

	accessToken, err := u.checkCredentials(c)
	if err != nil {
		ue.Err = err
		ue.Message = "wrong password or login"
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}

	u.AccessToken = accessToken
	err = u.upadteUserToDB(c)
	if err != nil {
		ue.Err = err
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}
