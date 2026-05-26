package models

import (
	"jamsel-backend/database"
)

type Order struct {
	ID              int64   `json:"id"`
	OrderNumber     string  `json:"order_number"`
	UserID          int64   `json:"user_id"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerPhone   string  `json:"customer_phone"`
	CustomerAddress string  `json:"customer_address"`
	CustomerCity    string  `json:"customer_city"`
	CustomerZip     string  `json:"customer_zip"`
	Subtotal        float64 `json:"subtotal"`
	Shipping        float64 `json:"shipping"`
	Total           float64 `json:"total"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

type OrderItem struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"order_id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	ProductImage string  `json:"product_image"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

// CreateOrder creates a new order and returns the order ID
func CreateOrder(order *Order) (int64, error) {
	query := `INSERT INTO orders (order_number, user_id, customer_name, customer_email, 
              customer_phone, customer_address, customer_city, customer_zip, 
              subtotal, shipping, total, payment_method, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
              RETURNING id`

	var orderID int64
	err := database.DB.QueryRow(query, order.OrderNumber, order.UserID,
		order.CustomerName, order.CustomerEmail, order.CustomerPhone,
		order.CustomerAddress, order.CustomerCity, order.CustomerZip,
		order.Subtotal, order.Shipping, order.Total, order.PaymentMethod, "pending").Scan(&orderID)

	return orderID, err
}

// AddOrderItem adds an item to an order
func AddOrderItem(orderItem *OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, product_name, 
              product_price, product_image, quantity, subtotal)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := database.DB.Exec(query, orderItem.OrderID, orderItem.ProductID,
		orderItem.ProductName, orderItem.ProductPrice, orderItem.ProductImage,
		orderItem.Quantity, orderItem.Subtotal)

	return err
}

// ClearCartAfterOrder removes all cart items for a user after order is placed
func ClearCartAfterOrder(userID int64) error {
	query := `DELETE FROM cart WHERE user_id = $1`
	_, err := database.DB.Exec(query, userID)
	return err
}
