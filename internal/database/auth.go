package database

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// InvalidCredentialsErr is returned from AuthUser if username or password is invalid
var InvalidCredentialsErr = errors.New("invalid credentials")

func (db Database) SignupUser(username, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := db.db.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	return err
}

// UserLogin will log in the user and return it's userID if it was ok
func (db Database) UserLogin(username, password string) (uint64, error) {
	var id uint64
	var hashedPassword string
	err := db.db.QueryRow(context.Background(), "SELECT id, password FROM users WHERE username=$1", username).Scan(&id, &hashedPassword)
	if err == pgx.ErrNoRows {
		return 0, InvalidCredentialsErr
	}
	if err != nil {
		return 0, errors.Wrap(err, "cannot query row")
	}
	// Check password
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return 0, InvalidCredentialsErr
	}
	return id, nil
}
