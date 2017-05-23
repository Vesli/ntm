package main

/*
	A userModel that represt the structu of a user.
	So far inserting a user into the DB is done with method.
	The methods will be refacto to have a more generic DB
	update matching the events routing
*/

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/vesli/ntm/helper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	ID          string    `json:"id" bson:"id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Password    string    `json:"password" bson:"password"`
	Email       string    `json:"email" bson:"email"`
	Permission  int8      `json:"permission" bson:"permission"`
	DateCreated time.Time `json:"date_creation" bson:"date_created"`
	AccessToken []byte    `json:"access_token,omitempty" bson:"access_token,omitempty"`
}

func (u *user) insertUserToDB(w http.ResponseWriter, c *mgo.Collection) {
	u.DateCreated = time.Now()

	h := md5.New()
	h.Write([]byte(u.Name + u.DateCreated.String()))
	u.ID = hex.EncodeToString(h.Sum(nil))

	index := mgo.Index{
		Key:        []string{u.ID},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		helper.WriteJSON(w, err, http.StatusInternalServerError)
	}

	err = c.Insert(u)
	if err != nil {
		helper.WriteJSON(w, err, http.StatusInternalServerError)
	}

	helper.WriteJSON(w, u, http.StatusOK)
}

func (u *user) updateUserToDB(c *mgo.Collection) error {
	err := c.Update(bson.M{"name": u.Name}, u)
	if err != nil {
		return err
	}

	return nil
}

func (u *user) findUserInDB(id string, w http.ResponseWriter, c *mgo.Collection) {
	err := c.Find(bson.M{"id": id}).One(u)
	if err != nil {
		helper.WriteJSON(w, "No user found", http.StatusInternalServerError)
		return
	}
	helper.WriteJSON(w, u, http.StatusOK)
}

func (u *user) userAlreadyExists(c *mgo.Collection) bool {
	err := c.Find(bson.M{"name": u.Name, "email": u.Email}).One(u)
	if err == nil {
		return true
	}

	return false
}

func (u *user) checkCredentials(c *mgo.Collection) ([]byte, error) {
	err := c.Find(bson.M{"name": u.Name, "password": u.Password}).One(u)
	if err != nil {
		return []byte(nil), errors.New("Wrong password or login")
	}

	h := sha256.New()
	h.Write([]byte(u.Password))
	return h.Sum(nil), nil
}
