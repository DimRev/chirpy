package database

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DimRev/chirpy/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserResp struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (db *DB) CreateUser(email, password string) (UserResp, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return UserResp{}, err
	}

	newUser := User{}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return UserResp{}, fmt.Errorf("failed hashing password: %v", err)
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
		return UserResp{}, err
	}

	token, err := auth.CreateToken(nil, newUser.Id)
	if err != nil {
		return UserResp{}, err
	}

	return UserResp{
		Id:    newUser.Id,
		Email: newUser.Email,
		Token: token,
	}, nil
}

func (db *DB) UpdateUser(email, password, tokenString string) (UserResp, error) {
	userIdStr, err := auth.ValidateJWT(tokenString)
	if err != nil {
		return UserResp{}, fmt.Errorf("error phrasing the session token: %v", err)
	}

	id, err := strconv.Atoi(userIdStr)
	if err != nil {
		return UserResp{}, err
	}

	dbContent, err := db.loadDB()
	if err != nil {
		return UserResp{}, err
	}

	prevUser, ok := dbContent.Users[id]
	if !ok {
		return UserResp{}, errors.New("User id not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return UserResp{}, fmt.Errorf("failed hashing password: %v", err)
	}

	dbContent.Users[id] = User{
		Id:               prevUser.Id,
		Password:         string(hashedPassword),
		Email:            email,
		ExpiresInSeconds: prevUser.ExpiresInSeconds,
	}

	err = db.writeDB(dbContent)
	if err != nil {
		return UserResp{}, err
	}

	return UserResp{
		Id:    id,
		Email: email,
		Token: tokenString,
	}, nil
}

func (db *DB) Login(email, password string, ExpiresInSeconds int) (UserResp, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return UserResp{}, err
	}

	for _, user := range dbContent.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return UserResp{}, errors.New("wrong email or password")
			}

			defaultExpiration := 60 * 60 * 24
			if ExpiresInSeconds == 0 {
				ExpiresInSeconds = defaultExpiration
			} else if ExpiresInSeconds > defaultExpiration {
				ExpiresInSeconds = defaultExpiration
			}

			token, err := auth.CreateToken(&ExpiresInSeconds, user.Id)
			if err != nil {
				return UserResp{}, errors.New("could not create JWT")
			}

			return UserResp{
				Id:    user.Id,
				Email: user.Email,
				Token: token,
			}, nil
		}
	}
	return UserResp{}, errors.New("wrong email or password")
}
