package handlers

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

var DB *sql.DB
var Store *sessions.CookieStore

func SetDB(db *sql.DB) {
	DB = db
}

func SetStore(store *sessions.CookieStore) {
	Store = store
}
