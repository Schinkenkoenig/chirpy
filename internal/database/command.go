package database

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) getUserByEmail(email string) (*User, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	for _, u := range db_structure.Users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, errors.New("user not found")
}

func (db *DB) IsPasswordCorrect(email string, password string) (*User, error) {
	u, err := db.getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) RevokeToken(userId int, token string) error {
	db_structure, err := db.loadDb()
	if err != nil {
		return err
	}

	u, ok := db_structure.Users[userId]

	if !ok {
		return errors.New("user not found")
	}

	u.RevokedTokens[token] = time.Now().UTC()

	err = db.writeDb(*db_structure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpgradeUser(userId int) (*User, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	// exist
	user, ok := db_structure.Users[userId]
	if !ok {
		return nil, errors.New("not found")
	}

	user.IsChirpyRed = true

	db_structure.Users[userId] = user

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) UpdateUser(userId int, email, password string) (*User, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	// exist
	user, ok := db_structure.Users[userId]
	if !ok {
		return nil, errors.New("not found")
	}

	hash, err := hashPassword(password)

	user.Email = email
	user.Password = hash

	if err != nil {
		return nil, err
	}

	u, _ := db.getUserByEmail(email)

	if u != nil {
		return nil, errors.New("user with this email already exists")
	}

	db_structure.Users[userId] = user

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) CreateUser(email string, password string) (*User, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	user_id := len(db_structure.Users) + 1
	hash, err := hashPassword(password)

	user := User{
		Email:         email,
		Id:            user_id,
		Password:      hash,
		RevokedTokens: make(map[string]time.Time),
	}

	if err != nil {
		return nil, err
	}

	u, _ := db.getUserByEmail(email)

	if u != nil {
		return nil, errors.New("user with this email already exists")
	}

	db_structure.Users[user_id] = user

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (db *DB) DeleteChirp(userId int) error {
	db_structure, err := db.loadDb()
	if err != nil {
		return err
	}

	delete(db_structure.Chirps, userId)

	err = db.writeDb(*db_structure)
	if err != nil {
		return err
	}

	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(userId int, body string) (*Chirp, error) {
	db_structure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirp_id := len(db_structure.Chirps) + 1
	chirp := Chirp{Body: body, Id: chirp_id, AuthorId: userId}

	db_structure.Chirps[chirp_id] = chirp

	err = db.writeDb(*db_structure)
	if err != nil {
		return nil, err
	}

	return &chirp, nil
}
