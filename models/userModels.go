package models

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

	FBID     string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	LoggedIn bool   `json:"logged_in"`

	Picture picture `json:"picture"`
	Events  []Event
}
