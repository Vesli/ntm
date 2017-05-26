package main

import "github.com/jinzhu/gorm"

// Event struct used as DB model
type Event struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	LocationX   float32 `json:"location_x"`
	LocationY   float32 `json:"location_y"`

	Users  []User `gorm:"many2many:event_users;"`
	User   User   `json:"created_by"`
	UserID int
}
