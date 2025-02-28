package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/Viridian-Software/budget_buddy/docs"
	"github.com/Viridian-Software/budget_buddy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

type apiConfig struct {
	conn_string string
	database    *database.Queries
	environment string
	jwtSecret   string
	port        string
}

// @title						Budget Buddy API
// @version					1.0
// @description				API Server for Budget Buddy Application
// @host						localhost:8080
// @BasePath					/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	err_loading_env := godotenv.Load()

	if err_loading_env != nil {
		log.Fatalf("error loading environment variables: %v", err_loading_env)
	}
	current_env := os.Getenv("ENVIRONMENT")
	var db_url string
	var secret string
	var server_port string
	if current_env == "dev" {
		db_url = os.Getenv("DEV_CONNECTION_STRING")
		secret = os.Getenv("JWTSECRET")
		server_port = os.Getenv("DEV_PORT")
	} else {
		db_url = os.Getenv("PROD_CONNECTION_STRING")
		secret = os.Getenv("JWTSECRET")
		server_port = os.Getenv("PROD_PORT")
	}

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

	// Swagger endpoint
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// User endpoints
	mux.HandleFunc("POST /api/users", config.AddUserHandler)
	mux.HandleFunc("PUT /api/users", config.UpdateUser)
	mux.HandleFunc("DELETE /api/users", config.DeleteUser)
	mux.HandleFunc("POST /api/login", config.UserLogin)
	mux.HandleFunc("POST /api/autoLogin", config.ValidateTokenHandler)
	mux.HandleFunc("POST /api/logout", config.LogoutHandler)
	mux.HandleFunc("GET /api/accounts/{userID}", config.GetAllUserAccounts)

	// Account endpoints
	mux.HandleFunc("POST /api/accounts", config.AddAccountHandler)
	mux.HandleFunc("DELETE /api/accounts", config.DeleteAccountHandler)
	mux.HandleFunc("POST /api/admin/reset", config.ResetUserTable)
	mux.HandleFunc("POST /api/refresh", config.HandleRefresh)
	mux.HandleFunc("POST /api/transactions", config.CreateTransaction)
	mux.HandleFunc("DELETE /api/transactions", config.DeleteTransaction)
	mux.HandleFunc("GET /api/transactions/{accountID}", config.GetTransactionsForAccount)

	// Admin endpoints
	mux.HandleFunc("POST /api/revoke", config.HandleRevoke)

	server := &http.Server{
		Addr:    ":" + config.port,
		Handler: mux,
	}
	server.ListenAndServe()
}
