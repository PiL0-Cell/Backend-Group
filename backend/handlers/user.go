package handlers

import (
	"encoding/json"
	"net/http"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not logged in"})
		return
	}

	session, _ := Store.Get(r, "jamsel-session")
	username, _ := session.Values["username"].(string)
	email, _ := session.Values["email"].(string)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"email":    email,
	})
}
