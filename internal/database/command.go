package database

func (db *DB) CreateUser(email string) (*User, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	user_id := len(db_structure.Users) + 1
	user := User{Email: email, Id: user_id}

	db_structure.Users[user_id] = user

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (*Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirp_id := len(db_structure.Chirps) + 1
	chirp := Chirp{Body: body, Id: chirp_id}

	db_structure.Chirps[chirp_id] = chirp

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &chirp, nil
}
