// package main

// import (
// 	"log"
// 	"net/http"

// 	"go-api/handlers"
// )

// func main() {

// 	mux := http.NewServeMux()

// 	// Routes
// 	mux.HandleFunc("/health", handlers.Health)

// 	// Users CRUD
// 	mux.HandleFunc("/users", handlers.Users)
// 	mux.HandleFunc("/users/", handlers.Users)

// 	log.Println("🚀 Server running on http://localhost:8080")

// 	log.Fatal(http.ListenAndServe(":8080", mux))
// }
//*************************************************************
// package main

// import (
// 	"log"
// 	"net/http"

// 	"go-api/config"
// 	"go-api/handlers"
// 	"go-api/models"

// 	"github.com/gorilla/mux"
// )

// func main() {

// 	// Connect DB
// 	config.ConnectDB()

// 	// Auto Create Table
// 	config.DB.AutoMigrate(&models.User{})

// 	// Router
// 	r := mux.NewRouter()

// 	// Routes
// 	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
// 	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
// 	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
// 	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
// 	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

// 	log.Println("🚀 Server running at http://localhost:8080")

// 	log.Fatal(http.ListenAndServe(":8080", r))
// }
//****************************************************************
package main

import (
	"log"
	"net/http"

	"go-api/config"
	"go-api/handlers"
	"go-api/middleware"
	"go-api/models"

	"github.com/gorilla/mux"
)

func main() {

	// DB
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// Router
	r := mux.NewRouter()

	// Logger (Global)
	r.Use(middleware.LoggerMiddleware)

	// Public
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// Protected
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	log.Println("🚀 Server running at :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}

