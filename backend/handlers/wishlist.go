package handlers

import (
	"encoding/json"
	"jamsel-backend/models"
	"jamsel-backend/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddToWishlist(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	var req struct {
		ProductID int `json:"product_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	err = models.AddToWishlist(userID, req.ProductID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add to wishlist"})
		return
	}

	go services.SyncUserWishlistToGorse(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Added to wishlist",
	})
}

func GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	wishlist, err := models.GetWishlist(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch wishlist"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wishlist)
}

func RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}

	err = models.RemoveFromWishlist(userID, productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove from wishlist"})
		return
	}

	go services.SyncUserWishlistToGorse(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Removed from wishlist",
	})
}
