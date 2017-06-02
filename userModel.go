package main

/*
	A userModel that represt the structu of a user.
	So far inserting a user into the DB is done with method.
	The methods will be refacto to have a more generic DB
	update matching the events routing
*/

import "github.com/jinzhu/gorm"

type data struct {
	URL string `json:"url"`
}

type picture struct {
	Data data `json:"Data"`
}

// User struct used as DB model
type User struct {
	gorm.Model
	FBID     string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Birthday string  `json:"birthday"`
	LoggedIn bool    `json:"logged_in"`
	Picture  picture `json:"picture"`
	Events   []Event
}

// userAlreadyExists check if a fb id is present.
// It's unique and sometimes other parameters can be blocked from user authorisation
func (u *User) userAlreadyExists(db *gorm.DB) bool {
	qs := db.First(&u, "fb_id = ?", u.FBID).GetErrors()
	if len(qs) == 0 {
		return true
	}

	return false
}

func (u *User) updateUserInDB(db *gorm.DB) error {
	err := db.Model(u).Update(u).Error
	if err != nil {
		return err
	}

	return nil
}
