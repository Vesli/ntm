package main

import "time"

type user struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Permission  int8
	DateCreated time.Time
}
