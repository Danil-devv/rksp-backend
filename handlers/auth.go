package handlers

import (
	"audio-call-auth-service/db"
	"audio-call-auth-service/models"
	"audio-call-auth-service/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s\n", err)
		}
	}(r.Body)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %s\n", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.DB.Exec(ctx, `
        INSERT INTO users (email, username, password_hash) 
        VALUES ($1, $2, $3)
    `, req.Email, req.Username, passwordHash)

	if err != nil {
		log.Printf("Error creating user: %s\n", err)
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s\n", err)
		}
	}(r.Body)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRow(ctx, `
        SELECT id, email, username, password_hash, created_at 
        FROM users WHERE username=$1
    `, req.Username).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		log.Printf("Invalid credentials: %s\n", err)
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		log.Printf("Invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	key, value, err := utils.GetSecretPair()
	if err != nil {
		http.Error(w, "Could not get secret pair", http.StatusInternalServerError)
		return
	}

	resp := SecretResponse{
		Key:   key,
		Value: value,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error writing login response: %s\n", err)
	}
}
