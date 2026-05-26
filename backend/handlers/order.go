package handlers

import (
	"encoding/json"
	"fmt"
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

	if req.OrderNumber == "" {
		req.OrderNumber = fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	}

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
		if err := models.AddOrderItem(orderItem); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add order items"})
			return
		}
	}

	if err := models.ClearCartAfterOrder(userID); err != nil {
		fmt.Println("Warning: Failed to clear cart:", err)
	}

	go services.SyncUserOrdersToGorse(userID)

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

	orders, err := models.GetUserOrders(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch orders"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
