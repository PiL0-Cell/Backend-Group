package models

import (
	"database/sql"
	db "jamsel-backend/database"
)

type Product struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Price         float64  `json:"price"`
	OriginalPrice *float64 `json:"original_price,omitempty"`
	Category      string   `json:"category"`
	Description   string   `json:"description"`
	Ingredients   string   `json:"ingredients"`
	Image         string   `json:"image"`
	Rating        float64  `json:"rating"`
	Reviews       int      `json:"reviews"`
	Badge         *string  `json:"badge,omitempty"`
	IsNew         bool     `json:"is_new"`
	IsSale        bool     `json:"is_sale"`
	Shades        string   `json:"shades"`
}

func GetAllProducts() ([]Product, error) {
	query := `SELECT id, name, price, original_price, category, description, 
              ingredients, image, rating, reviews, badge, is_new, is_sale, shades
              FROM products ORDER BY id`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.OriginalPrice, &p.Category,
			&p.Description, &p.Ingredients, &p.Image, &p.Rating, &p.Reviews,
			&p.Badge, &p.IsNew, &p.IsSale, &p.Shades)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(id int) (*Product, error) {
	query := `SELECT id, name, price, original_price, category, description, 
              ingredients, image, rating, reviews, badge, is_new, is_sale, shades
              FROM products WHERE id = $1`

	var p Product
	err := db.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.OriginalPrice, &p.Category,
		&p.Description, &p.Ingredients, &p.Image, &p.Rating, &p.Reviews,
		&p.Badge, &p.IsNew, &p.IsSale, &p.Shades)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
