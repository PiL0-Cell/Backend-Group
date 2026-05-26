package handlers

import (
	"encoding/json"
	"jamsel-backend/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	err = models.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add to cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Added to cart",
	})
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	cart, err := models.GetCart(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func UpdateCartQty(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Quantity int `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		models.RemoveFromCart(userID, productID)
	} else {
		err = models.UpdateCartQty(userID, productID, req.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update cart"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Cart updated",
	})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
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

	err = models.RemoveFromCart(userID, productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove from cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Removed from cart",
	})
}
