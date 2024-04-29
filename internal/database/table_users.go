package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, password string) (User, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUser := User{}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return User{}, fmt.Errorf("failed hashing password: %v", err)
	}

	for i := 1; ; i++ {
		_, ok := dbContent.Users[i]
		if !ok {
			newUser = User{
				Id:       i,
				Email:    email,
				Password: string(hashedPassword),
			}
			dbContent.Users[i] = newUser
			break
		}
	}

	err = db.writeDB(dbContent)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) Login(email string, password string) (User, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbContent.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return User{}, errors.New("wrong email or password")
			}
			return user, nil
		}
	}
	return User{}, errors.New("wrong email or password")
}
