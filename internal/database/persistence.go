package database

import (
	"encoding/json"
	"os"
)

func (db *DB) ensureDb() error {
	if _, err := os.Stat(db.path); os.IsExist(err) {
		return nil
	}
	db_structure := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	err := db.writeDb(db_structure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) writeDb(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

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

func (db *DB) loadDb() (*DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	var db_structure DBStructure

	data, err := os.ReadFile(db.path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &db_structure)
	if err != nil {
		return nil, err
	}

	return &db_structure, nil
}
