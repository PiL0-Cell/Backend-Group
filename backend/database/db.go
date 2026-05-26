package database

import (
	"database/sql"
)

var DB *sql.DB

func SetDB(database *sql.DB) {
	DB = database
}

func GetDB() *sql.DB {
	return DB
}
