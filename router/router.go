package router

import (
	// "go-postgres/router"
	"go-postgres/middleware"
	"net/http"

	"github.com/gorilla/mux"
)


// Router is exported and used in main.go
func Router() *mux.Router{
	router:= mux.NewRouter()

	router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock",middleware.GetAllStocks).Methods("GET","OPTIONS")
	router.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}",middleware.UpdateStock).Methods("PUT","OPTIONS")
	router.HandleFunc("/api/deletestock/{id}",middleware.DeleteStock).Methods("DELETE","OPTIONS")

	// ✅ Add a simple ping route to test router setup
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	// ✅ Add NotFound handler to catch undefined routes
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Route not found"))
	})
	return router
}