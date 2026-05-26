package models

import db "jamsel-backend/database"

type CartItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
}

func AddToCart(userID int64, productID int, quantity int) error {
	query := `INSERT INTO cart (user_id, product_id, quantity)
              VALUES ($1, $2, $3)
              ON CONFLICT (user_id, product_id)
              DO UPDATE SET quantity = cart.quantity + $3`

	_, err := db.DB.Exec(query, userID, productID, quantity)
	return err
}

func GetCart(userID int64) ([]CartItem, error) {
	query := `SELECT c.product_id, c.quantity, p.name, p.price, p.image
              FROM cart c
              JOIN products p ON c.product_id = p.id
              WHERE c.user_id = $1
              ORDER BY c.added_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity, &item.Name, &item.Price, &item.Image)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func UpdateCartQty(userID int64, productID int, quantity int) error {
	query := `UPDATE cart SET quantity = $3 WHERE user_id = $1 AND product_id = $2`
	_, err := db.DB.Exec(query, userID, productID, quantity)
	return err
}

func RemoveFromCart(userID int64, productID int) error {
	query := `DELETE FROM cart WHERE user_id = $1 AND product_id = $2`
	_, err := db.DB.Exec(query, userID, productID)
	return err
}
