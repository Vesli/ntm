package main

import (
	"encoding/json"
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

func registerUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var (
		u  user
		ue userException
	)

	err := decoder.Decode(&u)
	if err != nil {
		ue.Err = err
		ue.Message = "Error on decode"
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
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
	decoder := json.NewDecoder(r.Body)
	var (
		u  user
		ue userException
	)

	err := decoder.Decode(&u)
	if err != nil {
		ue.Err = err
		ue.Message = "Error on decode"
		helper.WriteJSON(w, ue, http.StatusBadRequest)
		return
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
