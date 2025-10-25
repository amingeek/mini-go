package main

import "time"

type User struct {
	FirstName   string    `form:"first_name" binding:"required,lendata"`
	LastName    string    `form:"last_name" binding:"required,lendata"`
	Username    string    `form:"username" binding:"required,alphanum,lenusername"`
	Email       string    `form:"email" binding:"required,email,lenemail"`
	PhoneNumber string    `form:"phone_number" binding:"required,number,startswith09,phone_number"`
	BirthDate   time.Time `form:"birth_date" binding:"required,birthdate" time_format:"2006/01/02"`
	NationalID  string    `form:"national_id" binding:"required,numeric,nationalid"`
}
