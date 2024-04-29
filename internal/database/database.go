package database

import (
	"encoding/json"
	"errors"
	"flag"
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
	Users  map[int]User  `json:"users"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	_, err := os.ReadFile(db.path)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) && errors.Is(pathError.Err, os.ErrNotExist) {
			mode := ""
			if *dbg {
				mode = "Debug"
			} else {
				mode = "Production"
			}
			log.Printf("%v mode: initializing new database", mode)
			db.writeDB(DBStructure{
				Chirps: make(map[int]Chirp),
				Users:  make(map[int]User),
			})
		} else {
			return fmt.Errorf("failed to read database file: %v", err)
		}
	} else if *dbg {
		db.writeDB(DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		})
		log.Println("Debug mode: wiping DB")
		log.Println("Debug mode: connecting to DB")
		return nil
	} else {
		log.Println("Production mode: connecting to DB")
	}

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
