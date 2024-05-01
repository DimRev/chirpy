package database

import (
	"errors"
	"fmt"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{}

	for i := 1; ; i++ {
		_, ok := dbContent.Chirps[i]
		if !ok {
			newChirp = Chirp{
				Id:       i,
				Body:     body,
				AuthorId: userId,
			}
			dbContent.Chirps[i] = newChirp
			break
		}
	}

	err = db.writeDB(dbContent)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, val := range dbContent.Chirps {
		chirps = append(chirps, val)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbContent.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("no chirp of id: %v in database", id)
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpId, userId int) error {
	dbContent, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbContent.Chirps[chirpId]
	if !ok {
		return errors.New("no such id found")
	}

	if chirp.AuthorId != userId {
		return errors.New("user not authorized")
	}

	newChirps := make(map[int]Chirp)
	for idx, currChirp := range dbContent.Chirps {
		if idx != chirpId {
			newChirps[idx] = currChirp
		}
	}

	dbContent.Chirps = newChirps

	return nil
}
