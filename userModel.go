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
	"time"

	"github.com/jinzhu/gorm"
)

type user struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Permission  int8      `json:"permission"`
	DateCreated time.Time `json:"date_creation"`
	AccessToken []byte    `json:"access_token,omitempty"`
}

func (u *user) userAlreadyExists(db *gorm.DB) bool {
	qs := db.First(&u, "email = ?", u.Email).GetErrors()
	if len(qs) == 0 {
		return true
	}

	return false
}

func (u *user) checkCredentials(db *gorm.DB) ([]byte, error) {
	err := db.First(&u, &u).Error
	if err != nil {
		return []byte(nil), errors.New("Wrong password or login")
	}

	h := sha256.New()
	h.Write([]byte(u.Password))
	return h.Sum(nil), nil
}
