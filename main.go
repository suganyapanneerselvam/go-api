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

	r.Use(middleware.LoggerMiddleware)

	// AUTH
	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/register", handlers.CreateUser).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// PROTECTED
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	handler := cors.AllowAll().Handler(r)

	log.Println("🚀 Server running at :8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
