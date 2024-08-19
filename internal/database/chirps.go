package database

import "fmt"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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
