package database

import (
	"sync"
	"time"
)

type DB struct {
	mux  *sync.RWMutex
	path string
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type User struct {
	RevokedTokens map[string]time.Time `json:"revoked_tokes"`
	Email         string               `json:"email"`
	Password      string               `json:"password"`
	Id            int                  `json:"id"`
}

type Chirp struct {
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
	Id       int    `json:"id"`
}

func NewDb(path string) (*DB, error) {
	// ensure db
	//
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
