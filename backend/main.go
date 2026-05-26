package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"jamsel-backend/handlers"
	"jamsel-backend/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var Store *sessions.CookieStore

func initDB() (*sql.DB, error) {
	// For Render production
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// If on Render, use SSL require
	if os.Getenv("RENDER") == "true" {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=prefer",
			host, port, user, password, dbname)
		log.Println("Connecting to Render PostgreSQL...")
		return sql.Open("postgres", connStr)
	}

	// Local development fallback
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}
	if dbname == "" {
		dbname = "jamsel_cosmetics"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", connStr)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://jamsel-backend.onrender.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic("Database ping failed: " + err.Error())
	}
	log.Println("Database connected successfully")

	// Initialize session store
	Store = sessions.NewCookieStore([]byte("jamsel-secret-key-change-this"))

	// Set global variables
	handlers.SetDB(db)
	handlers.SetStore(Store)

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/register", enableCORS(handlers.Register)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/login", enableCORS(handlers.Login)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/logout", enableCORS(handlers.Logout)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/products", enableCORS(handlers.GetAllProducts)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products/{id}", enableCORS(handlers.GetProductByID)).Methods("GET", "OPTIONS")

	// Protected routes
	router.HandleFunc("/api/user", enableCORS(handlers.GetCurrentUser)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cart", enableCORS(handlers.GetCart)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cart", enableCORS(handlers.AddToCart)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cart/{product_id}", enableCORS(handlers.UpdateCartQty)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/cart/{product_id}", enableCORS(handlers.RemoveFromCart)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/wishlist", enableCORS(handlers.GetWishlist)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/wishlist", enableCORS(handlers.AddToWishlist)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/wishlist/{product_id}", enableCORS(handlers.RemoveFromWishlist)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/orders", enableCORS(handlers.CreateOrder)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/orders", enableCORS(handlers.GetUserOrders)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cards", enableCORS(handlers.GetCreditCards)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cards", enableCORS(handlers.SaveCreditCard)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cards/default", enableCORS(handlers.GetDefaultCreditCard)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cards/{id}", enableCORS(handlers.DeleteCreditCard)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/cards/{id}/default", enableCORS(handlers.SetDefaultCreditCard)).Methods("PUT", "OPTIONS")

	// AI routes
	router.HandleFunc("/api/recommendations", enableCORS(handlers.GetAIRecommendations)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/sync-to-gorse", enableCORS(handlers.SyncUserToGorse)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/sync-products", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if err := services.SyncAllProducts(); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "Products synced to Gorse"})
	})).Methods("POST", "OPTIONS")

	// Serve frontend
	staticFile := http.FileServer(http.Dir("../frontend"))
	router.PathPrefix("/").Handler(staticFile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
