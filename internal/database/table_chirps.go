package database

import "fmt"

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{}

	for i := 1; ; i++ {
		_, ok := dbContent.Chirps[i]
		if !ok {
			newChirp = Chirp{
				Id:   i,
				Body: body,
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
