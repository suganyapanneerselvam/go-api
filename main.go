package main

import (
	"log"
	"net/http"

	"go-api/config"
	"go-api/handlers"
	"go-api/middleware"
	"go-api/models"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	config.ConnectDB()

	config.DB.AutoMigrate(&models.User{})

	r := mux.NewRouter()

	// Logger
	r.Use(middleware.LoggerMiddleware)

	// Public
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	handler := cors.AllowAll().Handler(r)

	log.Println("🚀 Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
