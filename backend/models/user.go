package models

import (
	"database/sql"
	"jamsel-backend/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) Create() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = database.DB.QueryRow(query, u.Username, u.Email, string(hashedPassword)).Scan(&u.ID)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	var u User
	var hashedPassword string

	err := database.DB.QueryRow(query, email).Scan(&u.ID, &u.Username, &u.Email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	u.Password = hashedPassword
	return &u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
