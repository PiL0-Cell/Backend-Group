package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"jamsel-backend/database"
	"jamsel-backend/models"
)

type GorseClient struct {
	BaseURL string
	Client  *http.Client
}

type Feedback struct {
	FeedbackType string    `json:"FeedbackType"`
	UserId       string    `json:"UserId"`
	ItemId       string    `json:"ItemId"`
	Value        float64   `json:"Value"`
	Timestamp    time.Time `json:"Timestamp"`
}

type Item struct {
	ItemId     string   `json:"ItemId"`
	IsHidden   bool     `json:"IsHidden,omitempty"`
	Labels     []string `json:"Labels,omitempty"`
	Categories []string `json:"Categories,omitempty"`
	Comment    string   `json:"Comment,omitempty"`
	Timestamp  string   `json:"Timestamp,omitempty"`
}

func NewGorseClient() *GorseClient {
	gorseURL := os.Getenv("GORSE_URL")
	if gorseURL == "" {
		// Local development fallback
		gorseURL = "http://localhost:8088"
	}
	log.Printf("Gorse client initialized with URL: %s", gorseURL)
	return &GorseClient{
		BaseURL: gorseURL,
		Client:  &http.Client{Timeout: 60 * time.Second}, // Increased timeout for Render cold starts
	}
}

func (g *GorseClient) InsertItem(item Item) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	url := g.BaseURL + "/api/item"
	log.Printf("Inserting item to Gorse: %s", url)

	resp, err := g.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Gorse InsertItem error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to insert item: %s", resp.Status)
	}
	return nil
}

func (g *GorseClient) InsertFeedback(feedback []Feedback) error {
	jsonData, err := json.Marshal(feedback)
	if err != nil {
		return err
	}

	url := g.BaseURL + "/api/feedback"
	log.Printf("Sending %d feedback items to Gorse", len(feedback))

	resp, err := g.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Gorse InsertFeedback error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Gorse error: %s", resp.Status)
	}
	return nil
}

func (g *GorseClient) GetRecommendations(userID string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/recommend/%s?n=%d", g.BaseURL, userID, n)

	resp, err := g.Client.Get(url)
	if err != nil {
		log.Printf("Gorse GetRecommendations error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return empty slice instead of error for no recommendations
		if resp.StatusCode == http.StatusNotFound {
			return []string{}, nil
		}
		return nil, fmt.Errorf("Gorse error: %s", resp.Status)
	}

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (g *GorseClient) HealthCheck() bool {
	url := g.BaseURL + "/health"
	resp, err := g.Client.Get(url)
	if err != nil {
		log.Printf("Gorse health check failed: %v", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func SyncAllProducts() error {
	gorse := NewGorseClient()

	// Check if Gorse is reachable
	if !gorse.HealthCheck() {
		log.Println("Gorse not reachable, skipping sync")
		return fmt.Errorf("Gorse not reachable at %s", gorse.BaseURL)
	}

	products, err := models.GetAllProducts()
	if err != nil {
		return err
	}

	log.Printf("Syncing %d products to Gorse...", len(products))

	successCount := 0
	for _, p := range products {
		item := Item{
			ItemId:     fmt.Sprintf("%d", p.ID),
			IsHidden:   false,
			Categories: []string{p.Category},
			Comment:    p.Name,
			Timestamp:  time.Now().Format(time.RFC3339),
		}
		if err := gorse.InsertItem(item); err != nil {
			log.Printf("Failed to sync product %d: %v", p.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully synced %d/%d products to Gorse", successCount, len(products))
	return nil
}

func SyncUserWishlistToGorse(userID int64) error {
	gorse := NewGorseClient()

	query := `SELECT product_id FROM wishlist WHERE user_id = $1`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var feedback []Feedback
	for rows.Next() {
		var productID int
		rows.Scan(&productID)
		feedback = append(feedback, Feedback{
			FeedbackType: "like",
			UserId:       fmt.Sprintf("%d", userID),
			ItemId:       fmt.Sprintf("%d", productID),
			Value:        1.0,
			Timestamp:    time.Now(),
		})
	}

	if len(feedback) > 0 {
		log.Printf("Syncing %d wishlist items for user %d", len(feedback), userID)
		return gorse.InsertFeedback(feedback)
	}
	return nil
}

func SyncUserOrdersToGorse(userID int64) error {
	gorse := NewGorseClient()

	query := `
        SELECT DISTINCT oi.product_id 
        FROM order_items oi
        JOIN orders o ON oi.order_id = o.id
        WHERE o.user_id = $1
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var feedback []Feedback
	for rows.Next() {
		var productID int
		rows.Scan(&productID)
		feedback = append(feedback, Feedback{
			FeedbackType: "purchase",
			UserId:       fmt.Sprintf("%d", userID),
			ItemId:       fmt.Sprintf("%d", productID),
			Value:        1.0,
			Timestamp:    time.Now(),
		})
	}

	if len(feedback) > 0 {
		log.Printf("Syncing %d purchased items for user %d", len(feedback), userID)
		return gorse.InsertFeedback(feedback)
	}
	return nil
}
