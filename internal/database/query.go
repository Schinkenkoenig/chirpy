package database

import (
	"cmp"
	"errors"
	"slices"
)

func (db *DB) GetChirpById(id int) (*Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	if v, ok := db_structure.Chirps[id]; ok {
		return &v, nil
	}

	return nil, errors.New("not found")
}

func (db *DB) IsTokenRevoked(userId int, token string) error {
	db_structure, err := db.loadDb()
	if err != nil {
		return err
	}

	if u, ok := db_structure.Users[userId]; ok {
		if _, ok := u.RevokedTokens[token]; ok {
			return errors.New("token revoked")
		}
		return nil
	}

	return errors.New("user not found")
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(db_structure.Chirps))

	for _, v := range db_structure.Chirps {
		chirps = append(chirps, v)
	}

	slices.SortFunc(chirps,
		func(a, b Chirp) int {
			return cmp.Compare(a.Id, b.Id)
		})

	return chirps, nil
}
