package database

import (
	"errors"
)

type User struct {
	Id               int    `json:"id"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
	IsChirpyRed      bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUser := User{}

	for i := 1; ; i++ {
		_, ok := dbContent.Users[i]
		if !ok {
			newUser = User{
				Id:          i,
				Email:       email,
				Password:    string(hashedPassword),
				IsChirpyRed: false,
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

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("User does not exist")
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, errors.New("User does not exist")
}

func (db *DB) UpdateUser(email, hashedPassword string, id int) (User, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	userToUpdate, err := db.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	updatedUser := User{
		Id:               id,
		Password:         hashedPassword,
		Email:            email,
		ExpiresInSeconds: userToUpdate.ExpiresInSeconds,
		IsChirpyRed:      userToUpdate.IsChirpyRed,
	}

	dbContent.Users[id] = updatedUser

	err = db.writeDB(dbContent)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) UpdateChirpyRed(id int) (User, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	userToUpdate, err := db.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	updatedUser := User{
		Id:               id,
		Email:            userToUpdate.Email,
		Password:         userToUpdate.Password,
		ExpiresInSeconds: userToUpdate.ExpiresInSeconds,
		IsChirpyRed:      true,
	}

	dbContent.Users[id] = updatedUser

	err = db.writeDB(dbContent)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}
