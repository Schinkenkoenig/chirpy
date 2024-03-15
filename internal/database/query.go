package database

func (db *DB) GetChirps() ([]Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(db_structure.Chirps))

	for _, v := range db_structure.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (*Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirp := Chirp{Body: body}

	db_structure.Chirps[len(db_structure.Chirps)] = chirp

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &chirp, nil
}
