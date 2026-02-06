package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	//"strconv"

	"go-api/config"
	"go-api/models"

	"github.com/gorilla/mux"
)

var jwtKey = []byte("my_secret_key")

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&creds)

	var user models.User

	// Find user
	if err := config.DB.Where("email = ?", creds.Email).
		First(&user).Error; err != nil {

		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(creds.Password),
	)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// Defaults
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{})

	// Search filter
	if search != "" {
		query = query.Where(
			"name ILIKE ? OR email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	// Count total
	query.Count(&total)

	// Fetch data
	query.
		Limit(limit).
		Offset(offset).
		Find(&users)

	// Response
	response := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  users,
	}

	json.NewEncoder(w).Encode(response)
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validation
	if strings.TrimSpace(user.Name) == "" ||
		strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" {
		http.Error(w, "Name, Email and Password are required", http.StatusBadRequest)
		return
	}

	// Email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Default role
	// Hash password
hashed, err := bcrypt.GenerateFromPassword(
	[]byte(user.Password),
	bcrypt.DefaultCost,
)
if err != nil {
	http.Error(w, "Password error", http.StatusInternalServerError)
	return
}

user.Password = string(hashed)
	user.Role = "user"

	// Save
	if err := config.DB.Create(&user).Error; err != nil {
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
