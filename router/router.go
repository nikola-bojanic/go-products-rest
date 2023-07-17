package router

import (
	"go-products-rest/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/products", middleware.GetProducts).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products/{id}", middleware.GetProduct).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products", middleware.CreateProduct).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/products/{id}", middleware.DeleteProduct).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/products/{id}", middleware.UpdateProduct).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/categories", middleware.GetCategories).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/categories/{id}", middleware.GetCategory).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/categories", middleware.CreateCategory).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/categories/{id}", middleware.DeleteCategory).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/categories/{id}", middleware.UpdateCategory).Methods("PUT", "OPTIONS")
	return router
}
