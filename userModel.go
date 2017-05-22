package main

/*
	A userModel that represt the structu of a user.
	So far inserting a user into the DB is done with method.
	The methods will be refacto to have a more generic DB
	update matching the events routing
*/

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/vesli/ntm/helper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Permission  int8      `json:"permission"`
	DateCreated time.Time `json:"date_creation"`
	AccessToken []byte
}

func (u *user) insertUserToDB(w http.ResponseWriter, c *mgo.Collection) {
	u.DateCreated = time.Now()
	err := c.Insert(u)
	if err != nil {
		helper.WriteJSON(w, err, http.StatusInternalServerError)
	}

	helper.WriteJSON(w, u, http.StatusOK)
}

func (u *user) upadteUserToDB(c *mgo.Collection) error {
	err := c.Update(bson.M{"name": u.Name}, u)
	if err != nil {
		return err
	}

	return nil
}

func (u *user) userAlreadyExists(c *mgo.Collection) bool {
	query := c.Find(bson.M{"name": u.Name, "email": u.Email}).One(u)
	if query == nil {
		return true
	}

	return false
}

func (u *user) checkCredentials(c *mgo.Collection) ([]byte, error) {
	query := c.Find(bson.M{"name": u.Name, "password": u.Password}).One(u)
	if query != nil {
		return []byte(nil), errors.New("Wrong password or login")
	}

	h := sha256.New()
	h.Write([]byte(u.Password))
	return h.Sum(nil), nil
}
