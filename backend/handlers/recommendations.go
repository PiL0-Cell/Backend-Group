package handlers

import (
	"encoding/json"
	"jamsel-backend/models"
	"jamsel-backend/services"
	"log"
	"net/http"
	"strconv"
)

func GetAIRecommendations(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAIRecommendations called")

	// Only return recommendations for logged-in users
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		log.Println("User not logged in, returning empty")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Product{})
		return
	}

	log.Printf("Getting recommendations for user: %d", userID)

	gorse := services.NewGorseClient()

	// Get personalized recommendations
	recommendedIDs, err := gorse.GetRecommendations(strconv.FormatInt(userID, 10), 8)
	if err != nil {
		log.Printf("Gorse GetRecommendations error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Product{})
		return
	}

	log.Printf("Received %d recommendations from Gorse", len(recommendedIDs))

	if len(recommendedIDs) == 0 {
		log.Println("No recommendations from Gorse")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Product{})
		return
	}

	// Convert IDs to product objects
	var recommendations []models.Product
	for _, idStr := range recommendedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		product, err := models.GetProductByID(id)
		if err == nil && product != nil {
			recommendations = append(recommendations, *product)
		}
		if len(recommendations) >= 6 {
			break
		}
	}

	log.Printf("Returning %d product recommendations", len(recommendations))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

func SyncUserToGorse(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not logged in"})
		return
	}

	log.Printf("Syncing user %d to Gorse", userID)

	go services.SyncUserWishlistToGorse(userID)
	go services.SyncUserOrdersToGorse(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "syncing"})
}
