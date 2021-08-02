package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const USERS_PATH = DB_PATH + "/users.json"
var ai AutoInc

type User struct {
	ID uint `json:"id"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetAllUsers() []User {
	f, err := os.OpenFile(USERS_PATH, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}
	defer f.Close()
	var users []User
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(data, &users); err != nil {
		return nil
	}
	return users
}

func GetUserByEmail(email string) *User {
	users := GetAllUsers()
	for _, u := range users {
		if u.Email == email {
			return &u
		}
	}
	return nil
}

func AddUser(u User) error {
	var users []User
	users = GetAllUsers()
	users = append(users, u)
	usersJson, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(USERS_PATH, usersJson, 0644)
}

func NewUser() *User {
	return &User {
		ID: ai.ID(),
	}
}