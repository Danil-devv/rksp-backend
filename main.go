package main

import (
	"audio-call-auth-service/db"
	"audio-call-auth-service/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Cannot load .env file: %s", err)
	}

	db.Connect(os.Getenv("PG_URL"))
	db.RunMigrations()

	r := mux.NewRouter()
	{
		r.HandleFunc("/register", handlers.Register).Methods(http.MethodPost, http.MethodOptions)
		r.HandleFunc("/login", handlers.Login).Methods(http.MethodPost, http.MethodOptions)
		r.Use(handlers.CorsMiddleware)
	}

	port := os.Getenv("PORT")
	log.Println("Server listening on port", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
