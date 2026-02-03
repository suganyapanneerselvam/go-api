// package handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"go-api/models"
// )

// // Fake Database
// var users = []models.User{
// 	{ID: 1, Name: "Alice", Email: "alice@test.com"},
// 	{ID: 2, Name: "Bob", Email: "bob@test.com"},
// }

// // Health Check
// func Health(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }

// // Main Users Handler
// func Users(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	// Get ID from URL (if exists)
// 	path := strings.TrimPrefix(r.URL.Path, "/users")
// 	path = strings.Trim(path, "/")

// 	// If ID exists → /users/1
// 	if path != "" {
// 		handleUserByID(w, r, path)
// 		return
// 	}

// 	// Otherwise → /users
// 	switch r.Method {

// 	// READ ALL
// 	case http.MethodGet:
// 		json.NewEncoder(w).Encode(users)

// 	// CREATE
// 	case http.MethodPost:

// 		var user models.User

// 		err := json.NewDecoder(r.Body).Decode(&user)
// 		if err != nil {
// 			http.Error(w, "Invalid body", http.StatusBadRequest)
// 			return
// 		}

// 		user.ID = len(users) + 1
// 		users = append(users, user)

// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(user)

// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// // Handle /users/{id}
// func handleUserByID(w http.ResponseWriter, r *http.Request, idStr string) {

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	index := -1

// 	// Find user index
// 	for i, u := range users {
// 		if u.ID == id {
// 			index = i
// 			break
// 		}
// 	}

// 	if index == -1 {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	switch r.Method {

// 	// READ ONE
// 	case http.MethodGet:
// 		json.NewEncoder(w).Encode(users[index])

// 	// UPDATE
// 	case http.MethodPut:

// 		var updatedUser models.User

// 		err := json.NewDecoder(r.Body).Decode(&updatedUser)
// 		if err != nil {
// 			http.Error(w, "Invalid body", http.StatusBadRequest)
// 			return
// 		}

// 		updatedUser.ID = id
// 		users[index] = updatedUser

// 		json.NewEncoder(w).Encode(updatedUser)

// 	// DELETE
// 	case http.MethodDelete:

// 		users = append(users[:index], users[index+1:]...)

// 		//w.WriteHeader(http.StatusNoContent)
// w.WriteHeader(http.StatusOK)

// json.NewEncoder(w).Encode(map[string]string{
// 	"message": "User deleted successfully",
// })

// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	//"strconv"

	"go-api/config"
	"go-api/models"

	"github.com/gorilla/mux"
)

// GET ALL USERS
func GetUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var users []models.User

	config.DB.Find(&users)

	json.NewEncoder(w).Encode(users)
}

// GET SINGLE USER
func GetUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id := params["id"]

	var user models.User

	result := config.DB.First(&user, id)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// CREATE USER
// func CreateUser(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	var user models.User

// 	err := json.NewDecoder(r.Body).Decode(&user)

// 	if err != nil {
// 		http.Error(w, "Invalid data", http.StatusBadRequest)
// 		return
// 	}

// 	config.DB.Create(&user)

// 	w.WriteHeader(http.StatusCreated)

// 	json.NewEncoder(w).Encode(user)
// }
func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var user models.User

	// Decode JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validation: Name
	if strings.TrimSpace(user.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validation: Email
	if strings.TrimSpace(user.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Email Format Check
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Save to DB
	result := config.DB.Create(&user)

	// Unique Email Error
	if result.Error != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UPDATE USER
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id := params["id"]

	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updated models.User

	json.NewDecoder(r.Body).Decode(&updated)

	user.Name = updated.Name
	user.Email = updated.Email

	config.DB.Save(&user)

	json.NewEncoder(w).Encode(user)
}

// DELETE USER
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id := params["id"]

	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	config.DB.Delete(&user)

	w.WriteHeader(http.StatusNoContent)
}
