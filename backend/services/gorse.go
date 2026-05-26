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
		gorseURL = "http://host.docker.internal:8088"
	}
	return &GorseClient{
		BaseURL: gorseURL,
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// HealthCheck tests if Gorse is running
func (g *GorseClient) HealthCheck() error {
	// Use /health, not /api/health
	resp, err := g.Client.Get(g.BaseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gorse health check failed: %s", resp.Status)
	}
	return nil
}

// InsertItem adds a product to Gorse's catalog
func (g *GorseClient) InsertItem(item Item) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	resp, err := g.Client.Post(
		g.BaseURL+"/api/item",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to insert item: %s", resp.Status)
	}
	return nil
}

// InsertFeedback sends user actions to Gorse
func (g *GorseClient) InsertFeedback(feedback []Feedback) error {
	jsonData, err := json.Marshal(feedback)
	if err != nil {
		return err
	}

	resp, err := g.Client.Post(
		g.BaseURL+"/api/feedback",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Gorse error: %s", resp.Status)
	}
	return nil
}

// GetRecommendations gets personalized recommendations for a user
func (g *GorseClient) GetRecommendations(userID string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/recommend/%s?n=%d", g.BaseURL, userID, n)

	resp, err := g.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gorse error: %s", resp.Status)
	}

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetLatestItems gets newest items
func (g *GorseClient) GetLatestItems(n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/latest?n=%d", g.BaseURL, n)

	resp, err := g.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gorse error: %s", resp.Status)
	}

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// SyncAllProducts syncs all products to Gorse
func SyncAllProducts() error {
	gorse := NewGorseClient()

	products, err := models.GetAllProducts()
	if err != nil {
		return err
	}

	log.Printf("Syncing %d products to Gorse...", len(products))

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
	}

	log.Println("Products synced to Gorse successfully")
	return nil
}

// SyncUserWishlistToGorse sends all wishlist items to Gorse
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
			FeedbackType: "like", // Gorse uses "like" for wishlist type actions
			UserId:       fmt.Sprintf("%d", userID),
			ItemId:       fmt.Sprintf("%d", productID),
			Value:        1.0,
			Timestamp:    time.Now(),
		})
	}

	if len(feedback) > 0 {
		log.Printf("Syncing %d wishlist items for user %d to Gorse", len(feedback), userID)
		return gorse.InsertFeedback(feedback)
	}
	return nil
}

// SyncUserOrdersToGorse sends all purchased items to Gorse
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
		log.Printf("Syncing %d purchased items for user %d to Gorse", len(feedback), userID)
		return gorse.InsertFeedback(feedback)
	}
	return nil
}
