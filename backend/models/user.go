package models

import (
	"database/sql"
	"jamsel-backend/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) Create() error {
	log.Printf("Hashing password for user: %s", u.Username)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("BCrypt error: %v", err)
		return err
	}

	log.Printf("Inserting user into database: %s, %s", u.Username, u.Email)

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = database.DB.QueryRow(query, u.Username, u.Email, string(hashedPassword)).Scan(&u.ID)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		return err
	}

	log.Printf("User created with ID: %d", u.ID)
	return nil
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
		log.Printf("GetUserByEmail error: %v", err)
		return nil, err
	}

	u.Password = hashedPassword
	return &u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
