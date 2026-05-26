package models

import (
	"crypto/sha256"
	"encoding/hex"
	"jamsel-backend/database"
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

// Hash card number for storage (never store raw card number)
func HashCardNumber(cardNumber string) string {
	hash := sha256.Sum256([]byte(cardNumber))
	return hex.EncodeToString(hash[:])
}

// Detect card brand from number
func DetectCardBrand(cardNumber string) string {
	if len(cardNumber) < 2 {
		return "Unknown"
	}
	firstDigit := cardNumber[0:1]
	firstTwo := cardNumber[0:2]
	firstFour := ""
	if len(cardNumber) >= 4 {
		firstFour = cardNumber[0:4]
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

// Save credit card
func SaveCreditCard(card *CreditCard) error {
	query := `INSERT INTO credit_cards (user_id, card_number_hash, card_last4, 
              card_holder_name, expiry_month, expiry_year, card_brand, is_default)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING id`

	err := database.DB.QueryRow(query, card.UserID, card.CardNumberHash, card.CardLast4,
		card.CardHolderName, card.ExpiryMonth, card.ExpiryYear, card.CardBrand, card.IsDefault).Scan(&card.ID)
	return err
}

// Get user's credit cards
func GetUserCreditCards(userID int64) ([]CreditCard, error) {
	query := `SELECT id, card_last4, card_holder_name, expiry_month, expiry_year, 
              card_brand, is_default, created_at
              FROM credit_cards WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []CreditCard
	for rows.Next() {
		var card CreditCard
		err := rows.Scan(&card.ID, &card.CardLast4, &card.CardHolderName,
			&card.ExpiryMonth, &card.ExpiryYear, &card.CardBrand, &card.IsDefault, &card.CreatedAt)
		if err != nil {
			return nil, err
		}
		card.UserID = userID
		cards = append(cards, card)
	}
	return cards, nil
}

// Get default credit card
func GetDefaultCreditCard(userID int64) (*CreditCard, error) {
	query := `SELECT id, card_last4, card_holder_name, expiry_month, expiry_year, 
              card_brand, is_default, created_at
              FROM credit_cards WHERE user_id = $1 AND is_default = true LIMIT 1`

	var card CreditCard
	err := database.DB.QueryRow(query, userID).Scan(&card.ID, &card.CardLast4, &card.CardHolderName,
		&card.ExpiryMonth, &card.ExpiryYear, &card.CardBrand, &card.IsDefault, &card.CreatedAt)
	if err != nil {
		return nil, err
	}
	card.UserID = userID
	return &card, nil
}

// Delete credit card
func DeleteCreditCard(cardID int64, userID int64) error {
	query := `DELETE FROM credit_cards WHERE id = $1 AND user_id = $2`
	_, err := database.DB.Exec(query, cardID, userID)
	return err
}

// Set card as default (and unset others)
func SetDefaultCreditCard(cardID int64, userID int64) error {
	// First, unset all defaults for this user
	_, err := database.DB.Exec(`UPDATE credit_cards SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Then set the selected card as default
	_, err = database.DB.Exec(`UPDATE credit_cards SET is_default = true WHERE id = $1 AND user_id = $2`, cardID, userID)
	return err
}
