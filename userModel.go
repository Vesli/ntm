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

	"github.com/jinzhu/gorm"
)

// User struct used as DB model
type User struct {
	gorm.Model
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Permission  int8   `json:"permission"`
	AccessToken []byte `json:"access_token,omitempty"`
	Events      []Event
}

func (u *User) userAlreadyExists(db *gorm.DB) bool {
	qs := db.First(&u, "email = ?", u.Email).GetErrors()
	if len(qs) == 0 {
		return true
	}

	return false
}

func (u *User) checkCredentials(db *gorm.DB) ([]byte, error) {
	err := db.First(&u, &u).Error
	if err != nil {
		return []byte(nil), errors.New("Wrong password or login")
	}

	h := sha256.New()
	h.Write([]byte(u.Password))

	return h.Sum(nil), nil
}

func (u *User) updateUserInDB(db *gorm.DB) error {
	err := db.Model(u).Update(u).Error
	if err != nil {
		return err
	}

	return nil
}
