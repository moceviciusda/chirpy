package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithoutPassword struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (db *DB) CreateUser(email, password string) (UserWithoutPassword, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return UserWithoutPassword{}, err
	}

	for _, usr := range dbStruct.Users {
		if usr.Email == email {
			return UserWithoutPassword{}, fmt.Errorf("email address: %s is already in use", email)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{id, email, string(hash)}

	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	return UserWithoutPassword{user.Id, user.Email, ""}, nil
}

// func (db *DB) UpdateUser(email, password string) (UserWithoutPassword, error) {

// }

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, usr := range dbStruct.Users {
		if usr.Email == email {
			return usr, nil
		}
	}

	return User{}, fmt.Errorf("user does not exist")
}
