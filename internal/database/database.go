package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	mux := &sync.RWMutex{}
	db := DB{
		path: path,
		mux:  mux,
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()

	_, err := os.ReadFile(db.path)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) && errors.Is(pathError.Err, os.ErrNotExist) {
			log.Printf("creating new db-file: %v", db.path)
			emptyChirps := DBStructure{
				Chirps: make(map[int]Chirp),
			}
			jsonData, err := json.Marshal(emptyChirps)
			if err != nil {
				return fmt.Errorf("failed to create new database file: %v", err)
			}
			os.WriteFile(db.path, []byte(jsonData), 0755)
		} else {
			return fmt.Errorf("failed to read database file: %v", err)
		}
	}

	db.mux.Unlock()
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()

	jsonContent, err := os.ReadFile(db.path)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) && errors.Is(pathError.Err, os.ErrNotExist) {
			return DBStructure{}, fmt.Errorf("failed to find database file: %v", err)
		}
		return DBStructure{}, fmt.Errorf("failed to read database file: %v", err)
	}
	dbContent := DBStructure{}
	err = json.Unmarshal(jsonContent, &dbContent)
	if err != nil {
		return DBStructure{}, fmt.Errorf("\nfailed to Unmarshal json content: %v", err)
	}

	db.mux.RUnlock()
	return dbContent, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("\nfailed to Marshal json content: %v", err)
	}

	db.mux.Lock()
	err = os.WriteFile(db.path, dat, 0755)
	if err != nil {
		return fmt.Errorf("\nfailed to write to file: %v", err)
	}
	db.mux.Unlock()
	return nil
}
