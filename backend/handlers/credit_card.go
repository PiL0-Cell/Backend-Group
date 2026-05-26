package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"jamsel-backend/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CreditCard struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	CardNumberHash string `json:"-"`
	CardLast4      string `json:"card_last4"`
	CardHolderName string `json:"card_holder_name"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CardBrand      string `json:"card_brand"`
	IsDefault      bool   `json:"is_default"`
	CreatedAt      string `json:"created_at"`
}

type SaveCardRequest struct {
	CardNumber     string `json:"card_number"`
	CardHolderName string `json:"card_holder_name"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	IsDefault      bool   `json:"is_default"`
}

func hashCardNumber(cardNumber string) string {
	hash := sha256.Sum256([]byte(cardNumber))
	return hex.EncodeToString(hash[:])
}

func detectCardBrand(cardNumber string) string {
	clean := strings.ReplaceAll(cardNumber, " ", "")
	if len(clean) < 2 {
		return "Unknown"
	}
	firstDigit := clean[0:1]
	firstTwo := clean[0:2]
	firstFour := ""
	if len(clean) >= 4 {
		firstFour = clean[0:4]
	}

	switch {
	case firstDigit == "4":
		return "Visa"
	case firstTwo >= "51" && firstTwo <= "55":
		return "Mastercard"
	case firstTwo == "34" || firstTwo == "37":
		return "American Express"
	case firstFour == "6011" || firstTwo == "65":
		return "Discover"
	default:
		return "Unknown"
	}
}

func isValidLuhn(cardNumber string) bool {
	clean := strings.ReplaceAll(cardNumber, " ", "")
	var sum int
	var alternate bool

	for i := len(clean) - 1; i >= 0; i-- {
		n := int(clean[i] - '0')
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

	cardNumber := strings.ReplaceAll(req.CardNumber, " ", "")
	if !isValidLuhn(cardNumber) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid card number"})
		return
	}

	last4 := cardNumber[len(cardNumber)-4:]
	expiryMonth, _ := strconv.Atoi(req.ExpiryMonth)
	expiryYear, _ := strconv.Atoi(req.ExpiryYear)
	brand := detectCardBrand(cardNumber)

	query := `INSERT INTO credit_cards (user_id, card_number_hash, card_last4, 
              card_holder_name, expiry_month, expiry_year, card_brand, is_default)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING id`

	var cardID int64
	err = database.DB.QueryRow(query, userID, hashCardNumber(cardNumber), last4,
		req.CardHolderName, expiryMonth, expiryYear, brand, req.IsDefault).Scan(&cardID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save card"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"card_id": cardID,
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

	query := `SELECT id, card_last4, card_holder_name, expiry_month, expiry_year, 
              card_brand, is_default, created_at
              FROM credit_cards WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cards"})
		return
	}
	defer rows.Close()

	var cards []CreditCard
	for rows.Next() {
		var card CreditCard
		err := rows.Scan(&card.ID, &card.CardLast4, &card.CardHolderName,
			&card.ExpiryMonth, &card.ExpiryYear, &card.CardBrand, &card.IsDefault, &card.CreatedAt)
		if err != nil {
			continue
		}
		card.UserID = userID
		cards = append(cards, card)
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

	query := `SELECT id, card_last4, card_holder_name, expiry_month, expiry_year, 
              card_brand, is_default, created_at
              FROM credit_cards WHERE user_id = $1 AND is_default = true LIMIT 1`

	var card CreditCard
	err := database.DB.QueryRow(query, userID).Scan(&card.ID, &card.CardLast4, &card.CardHolderName,
		&card.ExpiryMonth, &card.ExpiryYear, &card.CardBrand, &card.IsDefault, &card.CreatedAt)
	if err != nil {
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

	query := `DELETE FROM credit_cards WHERE id = $1 AND user_id = $2`
	_, err = database.DB.Exec(query, cardID, userID)
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

	_, err = database.DB.Exec(`UPDATE credit_cards SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update"})
		return
	}

	_, err = database.DB.Exec(`UPDATE credit_cards SET is_default = true WHERE id = $1 AND user_id = $2`, cardID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to set default"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Default card updated",
	})
}
