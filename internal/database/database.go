package database

import (
	"sync"
)

type DB struct {
	mux  *sync.RWMutex
	path string
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

func NewDb(path string) (*DB, error) {
	// ensure db
	db := DB{
		mux:  &sync.RWMutex{},
		path: path,
	}

	err := db.ensureDb()
	if err != nil {
		return nil, err
	}

	return &db, nil
}
