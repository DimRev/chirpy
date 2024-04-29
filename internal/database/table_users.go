package database

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	token, err := db.createToken(nil, newUser.Id)
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
	token, err := db.parseToken(tokenString)
	if err != nil {
		return UserResp{}, fmt.Errorf("error phrasing the session token: %v", err)
	}

	userIdStr := token.Subject
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

	return UserResp{
		Id:    id,
		Email: email,
		Token: tokenString,
	}, nil
}

func (db *DB) Login(email, password string) (UserResp, error) {
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

			token, err := db.createToken(user.ExpiresInSeconds, user.Id)
			if err != nil {
				return UserResp{}, err
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

func (db *DB) createToken(ExpiresInSeconds *int, id int) (string, error) {
	timeToAdd := 24 * time.Hour
	if ExpiresInSeconds != nil {
		timeToAdd = time.Duration(*ExpiresInSeconds) * time.Second
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeToAdd)),
		Subject:   fmt.Sprint(id),
	})

	signedToken, err := token.SignedString([]byte(db.jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (db *DB) parseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	// Parse the token. Second argument is a function that takes a parsed token
	// and returns the key for validating the token's signature.
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Check if the token's signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(db.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
