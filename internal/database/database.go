package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path  string
	mutex *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	mutex := sync.RWMutex{}
	db := DB{path, &mutex}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{id, body}

	dbStruct.Chirps[id] = chirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return chirp, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		chirps[chirp.Id-1] = chirp
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		return chirp, fmt.Errorf("chirp does not exist. ID: %v", id)
	}

	return chirp, nil
}

func (db *DB) ensureDB() error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		log.Println("DB not found. Initializing DB")
		contents := DBStructure{map[int]Chirp{}}
		data, err := json.Marshal(&contents)
		if err != nil {
			return err
		}

		err = os.WriteFile(db.path, data, 0666)
		if err != nil {
			return err
		}
	}

	log.Println("Connected to DB")

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	dbStruct := DBStructure{}

	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return dbStruct, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
