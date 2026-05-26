package handlers

import (
	"encoding/json"
	"jamsel-backend/models"
	"jamsel-backend/services"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetLoggedInUserID returns user ID from session cookie
func GetLoggedInUserID(r *http.Request) int64 {
	if Store == nil {
		log.Println("Store is nil in GetLoggedInUserID")
		return 0
	}

	session, err := Store.Get(r, "jamsel-session")
	if err != nil {
		log.Println("Session error:", err)
		return 0
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return 0
	}
	return userID
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
		return
	}

	if len(req.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password must be at least 6 characters"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err = user.Create()
	if err != nil {
		log.Println("Registration error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user. Email may already exist."})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user_id": user.ID,
		"message": "Account created successfully",
	})
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

	// Create session
	if Store == nil {
		log.Println("Store is nil in Login")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Session store not initialized"})
		return
	}

	session, _ := Store.Get(r, "jamsel-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	err = session.Save(r, w)
	if err != nil {
		log.Println("Session save error:", err)
	}

	go func() {
		if err := services.SyncUserWishlistToGorse(user.ID); err != nil {
			log.Println("Failed to sync wishlist to Gorse:", err)
		}
		if err := services.SyncUserOrdersToGorse(user.ID); err != nil {
			log.Println("Failed to sync orders to Gorse:", err)
		}
	}()

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
