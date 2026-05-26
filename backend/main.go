package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"jamsel-backend/database"
	"jamsel-backend/handlers"
	"jamsel-backend/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var Store *sessions.CookieStore

func init() {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=jamsel_cosmetics sslmode=disable"

	databaseConn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = databaseConn.Ping(); err != nil {
		panic("Database ping failed: " + err.Error())
	}

	database.SetDB(databaseConn)
	log.Println("Database connected successfully")
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow credentials (cookies)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
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
	// Initialize session store
	Store = sessions.NewCookieStore([]byte("jamsel-secret-key-change-this"))

	// Set store and DB in handlers
	handlers.SetStore(Store)
	handlers.SetDB(database.GetDB())

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/register", enableCORS(handlers.Register)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/login", enableCORS(handlers.Login)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/logout", enableCORS(handlers.Logout)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/products", enableCORS(handlers.GetAllProducts)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products/{id}", enableCORS(handlers.GetProductByID)).Methods("GET", "OPTIONS")

	// Protected routes (require login)
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
	// AI Recommendations routes
	router.HandleFunc("/api/recommendations", enableCORS(handlers.GetAIRecommendations)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/sync-to-gorse", enableCORS(handlers.SyncUserToGorse)).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/sync-products", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		err := services.SyncAllProducts()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "Products synced to Gorse"})
	})).Methods("POST", "OPTIONS")

	// Serve frontend
	staticFile := http.FileServer(http.Dir("../frontend"))
	router.PathPrefix("/").Handler(staticFile)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
