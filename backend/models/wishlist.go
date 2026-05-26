package models

import db "jamsel-backend/database"

type WishlistItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
}

func AddToWishlist(userID int64, productID int) error {
	query := `INSERT INTO wishlist (user_id, product_id)
              VALUES ($1, $2)
              ON CONFLICT (user_id, product_id) DO NOTHING`

	_, err := db.DB.Exec(query, userID, productID)
	return err
}

func GetWishlist(userID int64) ([]WishlistItem, error) {
	query := `SELECT w.product_id, p.name, p.price, p.image
              FROM wishlist w
              JOIN products p ON w.product_id = p.id
              WHERE w.user_id = $1
              ORDER BY w.added_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []WishlistItem
	for rows.Next() {
		var item WishlistItem
		err := rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Image)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func RemoveFromWishlist(userID int64, productID int) error {
	query := `DELETE FROM wishlist WHERE user_id = $1 AND product_id = $2`
	_, err := db.DB.Exec(query, userID, productID)
	return err
}
