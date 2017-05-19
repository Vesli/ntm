package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vesli/ntm/config"
	"github.com/vesli/ntm/helper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userException struct {
	Message string
}

const userCollection = "users"

func contextFromMiddleware(r *http.Request) (*config.Config, *mgo.Session) {
	conf := r.Context().Value(configuration).(*config.Config)
	sessionC := r.Context().Value(mgoSession).(*mgo.Session)

	return conf, sessionC
}

func userAlreadyExists(u *user, c *mgo.Collection) bool {
	query := c.Find(bson.M{"name": u.Name, "email": u.Email}).One(u)
	if query == nil {
		return true
	}

	return false
}

func insertUserToDB(w http.ResponseWriter, u user, c *mgo.Collection) {
	u.DateCreated = time.Now()
	err := c.Insert(u)
	if err != nil {
		helper.WriteJSON(w, err, 500)
	}

	helper.WriteJSON(w, u, 200)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var (
		u  user
		ue userException
	)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&u)
	if err != nil {
		helper.WriteJSON(w, err, 400)
		return
	}

	conf, sessionC := contextFromMiddleware(r)
	c := sessionC.DB(conf.DBName).C(userCollection)

	if userAlreadyExists(&u, c) {
		ue.Message = "User name or email already exists"
		helper.WriteJSON(w, ue, 400)
		return
	}

	insertUserToDB(w, u, c)
}
