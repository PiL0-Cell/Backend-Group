package handlers

import (
	"encoding/json"
	"jamsel-backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type SaveCardRequest struct {
	CardNumber     string `json:"card_number"`
	CardHolderName string `json:"card_holder_name"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	IsDefault      bool   `json:"is_default"`
}

func SaveCreditCard(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	var req SaveCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Remove spaces from card number
	cardNumber := strings.ReplaceAll(req.CardNumber, " ", "")

	// Validate card number (Luhn algorithm)
	if !isValidLuhn(cardNumber) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid card number"})
		return
	}

	// Get last 4 digits
	last4 := cardNumber[len(cardNumber)-4:]

	// Parse expiry
	expiryMonth, _ := strconv.Atoi(req.ExpiryMonth)
	expiryYear, _ := strconv.Atoi(req.ExpiryYear)

	// Validate expiry date
	if expiryMonth < 1 || expiryMonth > 12 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid expiry month"})
		return
	}

	// Detect card brand
	brand := models.DetectCardBrand(cardNumber)

	card := &models.CreditCard{
		UserID:         userID,
		CardNumberHash: models.HashCardNumber(cardNumber),
		CardLast4:      last4,
		CardHolderName: req.CardHolderName,
		ExpiryMonth:    expiryMonth,
		ExpiryYear:     expiryYear,
		CardBrand:      brand,
		IsDefault:      req.IsDefault,
	}

	err = models.SaveCreditCard(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save card"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"card_id": card.ID,
		"message": "Card saved successfully",
	})
}

func GetCreditCards(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	cards, err := models.GetUserCreditCards(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cards"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

func GetDefaultCreditCard(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	card, err := models.GetDefaultCreditCard(userID)
	if err != nil {
		// No default card is fine
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(card)
}

func DeleteCreditCard(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	vars := mux.Vars(r)
	cardID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid card ID"})
		return
	}

	err = models.DeleteCreditCard(cardID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete card"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Card deleted successfully",
	})
}

func SetDefaultCreditCard(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	vars := mux.Vars(r)
	cardID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid card ID"})
		return
	}

	err = models.SetDefaultCreditCard(cardID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to set default card"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Default card updated",
	})
}

// Luhn algorithm for credit card validation
func isValidLuhn(cardNumber string) bool {
	var sum int
	var alternate bool

	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')
		if alternate {
			n *= 2
			if n > 9 {
				n = n%10 + 1
			}
		}
		sum += n
		alternate = !alternate
	}
	return sum%10 == 0
}
