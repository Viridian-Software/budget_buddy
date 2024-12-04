package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Viridian-Software/budget_buddy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	conn_string string
	database    *database.Queries
	environment string
	jwtSecret   string
}

func main() {
	err_loading_env := godotenv.Load()

	if err_loading_env != nil {
		log.Fatalf("error loading environment variables: %v", err_loading_env)
	}

	db_url := os.Getenv("CONNECTION_STRING")
	current_env := os.Getenv("ENVIRONMENT")
	secret := os.Getenv("JWTSECRET")

	db, err_opening_db := sql.Open("postgres", db_url)
	if err_opening_db != nil {
		log.Fatalf("error opening database connection: %v", err_opening_db)
	}
	dbQueries := database.New(db)

	config := &apiConfig{
		conn_string: db_url,
		database:    dbQueries,
		environment: current_env,
		jwtSecret:   secret,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("/", ServerRunningHandler)
	mux.HandleFunc("POST /api/users", config.AddUserHandler)
	mux.HandleFunc("POST /api/login", config.UserLogin)
	server.ListenAndServe()
}
