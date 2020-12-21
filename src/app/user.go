package main

import (
	"log"
	"database/sql"
)

type User struct {
	User_id int `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Password string `json:"password,omitempty"`
	Email string `json:"email,omitempty"`
}

func (u *User) populate(rows *sql.Rows) {
	err := rows.Scan(&u.User_id, &u.FirstName, &u.LastName, &u.Password, &u.Email)
	if err != nil {
		log.Print(err)
	}
}