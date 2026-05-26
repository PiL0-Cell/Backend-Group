package handlers

import (
	"encoding/json"
	"jamsel-backend/database"
	"jamsel-backend/models"
	"jamsel-backend/services"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// Add panic recovery
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("PANIC in Register: %v", rec)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
	}()

	log.Println("Register endpoint called")

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON decode error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	log.Printf("Register request: username=%s, email=%s", req.Username, req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Println("Missing fields")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
		return
	}

	if len(req.Password) < 6 {
		log.Println("Password too short")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password must be at least 6 characters"})
		return
	}

	// Check if database is connected
	if database.DB == nil {
		log.Println("Database connection is nil!")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database not connected"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err = user.Create()
	if err != nil {
		log.Println("User.Create error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user. Email may already exist."})
		return
	}

	log.Printf("User created successfully: ID=%d", user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user_id": user.ID,
		"message": "Account created successfully",
	})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetLoggedInUserID(r *http.Request) int64 {
	session, err := Store.Get(r, "jamsel-session")
	if err != nil {
		return 0
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return 0
	}
	return userID
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email and password are required"})
		return
	}

	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		log.Println("Login error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	if user == nil || !user.CheckPassword(req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	session, _ := Store.Get(r, "jamsel-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Save(r, w)

	// Sync user data to Gorse
	go services.SyncUserWishlistToGorse(user.ID)
	go services.SyncUserOrdersToGorse(user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"message":  "Login successful",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "jamsel-session")
	session.Values = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}
