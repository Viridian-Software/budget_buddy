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
	port        string
}

func main() {
	err_loading_env := godotenv.Load()

	if err_loading_env != nil {
		log.Fatalf("error loading environment variables: %v", err_loading_env)
	}

	db_url := os.Getenv("CONNECTION_STRING")
	current_env := os.Getenv("ENVIRONMENT")
	secret := os.Getenv("JWTSECRET")
	server_port := os.Getenv("PORT")

	db, err_opening_db := sql.Open("postgres", db_url)
	if err_opening_db != nil {
		log.Fatalf("error opening database connection: %v", err_opening_db)
	}
	dbQueries := database.New(db)

	config := apiConfig{
		conn_string: db_url,
		database:    dbQueries,
		environment: current_env,
		jwtSecret:   secret,
		port:        server_port,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", config.AddUserHandler)
	mux.HandleFunc("PUT /api/users", config.UpdateUser)
	mux.HandleFunc("DELETE /api/users", config.DeleteUser)
	mux.HandleFunc("POST /api/login", config.UserLogin)
	mux.HandleFunc("POST /api/accounts", config.AddAccountHandler)
	mux.HandleFunc("DELETE /api/accounts", config.AddAccountHandler)
	mux.HandleFunc("POST /api/admin/reset", config.ResetUserTable)
	mux.HandleFunc("POST /api/refresh", config.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", config.HandleRevoke)
	mux.HandleFunc("GET /api/accounts/{userID}", config.GetAllUserAccounts)
	mux.HandleFunc("POST /api/transactions", config.CreateTransaction)
	mux.HandleFunc("DELETE /api/transactions", config.DeleteTransaction)
	server := &http.Server{
		Addr:    ":" + config.port,
		Handler: mux,
	}
	server.ListenAndServe()
}
