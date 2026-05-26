package handlers

import (
	"encoding/json"
	"fmt"
	"jamsel-backend/database"
	"jamsel-backend/models"
	"jamsel-backend/services"
	"net/http"
	"time"
)

type OrderRequest struct {
	OrderNumber string `json:"order_number"`
	Customer    struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
	} `json:"customer"`
	Items []struct {
		ProductID int     `json:"product_id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Quantity  int     `json:"quantity"`
		Image     string  `json:"image"`
	} `json:"items"`
	Subtotal float64 `json:"subtotal"`
	Shipping float64 `json:"shipping"`
	Total    float64 `json:"total"`
	Payment  string  `json:"payment"`
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	var req OrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Create order number if not provided
	if req.OrderNumber == "" {
		req.OrderNumber = fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	}

	// Create order in database
	order := &models.Order{
		OrderNumber:     req.OrderNumber,
		UserID:          userID,
		CustomerName:    req.Customer.Name,
		CustomerEmail:   req.Customer.Email,
		CustomerPhone:   req.Customer.Phone,
		CustomerAddress: req.Customer.Address,
		CustomerCity:    req.Customer.City,
		CustomerZip:     req.Customer.Zip,
		Subtotal:        req.Subtotal,
		Shipping:        req.Shipping,
		Total:           req.Total,
		PaymentMethod:   req.Payment,
		Status:          "pending",
	}

	orderID, err := models.CreateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create order"})
		return
	}

	// Add each item to order_items
	for _, item := range req.Items {
		orderItem := &models.OrderItem{
			OrderID:      orderID,
			ProductID:    item.ProductID,
			ProductName:  item.Name,
			ProductPrice: item.Price,
			ProductImage: item.Image,
			Quantity:     item.Quantity,
			Subtotal:     item.Price * float64(item.Quantity),
		}

		err = models.AddOrderItem(orderItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add order items"})
			return
		}
	}

	// Clear the user's cart after successful order
	err = models.ClearCartAfterOrder(userID)
	if err != nil {
		// Log error but don't fail the order
		fmt.Println("Warning: Failed to clear cart:", err)
	}

	go services.SyncUserWishlistToGorse(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"order_id":     orderID,
		"order_number": req.OrderNumber,
		"message":      "Order placed successfully",
	})
}

func GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedInUserID(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	query := `SELECT id, order_number, total, status, created_at 
              FROM orders WHERE user_id = $1 
              ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id int64
		var orderNumber string
		var total float64
		var status string
		var createdAt string

		rows.Scan(&id, &orderNumber, &total, &status, &createdAt)

		orders = append(orders, map[string]interface{}{
			"id":           id,
			"order_number": orderNumber,
			"total":        total,
			"status":       status,
			"created_at":   createdAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
