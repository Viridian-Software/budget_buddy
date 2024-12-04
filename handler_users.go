package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Viridian-Software/budget_buddy/internal/auth"
	"github.com/Viridian-Software/budget_buddy/internal/custom_errors"
	"github.com/Viridian-Software/budget_buddy/internal/database"
	"github.com/Viridian-Software/budget_buddy/internal/validators"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Is_Admin   bool      `json:"is_admin"`
	First_Name string    `json:"first_name"`
	Last_Name  string    `json:"last_name"`
}

type AddUser struct {
	Email      string `json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Password   string `json:"password"`
}

func (cfg *apiConfig) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	new_user_info := &AddUser{}
	decoder := json.NewDecoder(r.Body)
	err_decoding_response := decoder.Decode(&new_user_info)
	if err_decoding_response != nil {
		custom_errors.ReturnErrorWithMessage(w, "error decoding json", err_decoding_response, 500)
		return
	}
	if !validators.ValidatePassword(new_user_info.Password) {
		custom_errors.ReturnErrorWithMessage(w, "password does not meet requirements", nil, 401)
	}
	hashed_password, err_hashing_password := auth.HashPassword(new_user_info.Password)
	if err_hashing_password != nil {
		custom_errors.ReturnErrorWithMessage(w, "error hashing password", err_hashing_password, 500)
		return
	}
	dbUser, err_adding_user := cfg.database.AddUser(r.Context(), database.AddUserParams{
		Email:          new_user_info.Email,
		FirstName:      new_user_info.First_Name,
		LastName:       new_user_info.Last_Name,
		HashedPassword: hashed_password,
	})
	if err_adding_user != nil {
		custom_errors.ReturnErrorWithMessage(w, "error adding user", err_adding_user, 500)
		return
	}
	added_user := &User{
		ID:         dbUser.ID,
		Email:      dbUser.Email,
		First_Name: dbUser.FirstName,
		Last_Name:  dbUser.LastName,
		Created_At: dbUser.CreatedAt,
		Updated_At: dbUser.UpdatedAt,
		Is_Admin:   dbUser.IsAdmin,
	}
	jsonData, err_marshalling_json := json.Marshal(added_user)
	if err_marshalling_json != nil {
		custom_errors.ReturnErrorWithMessage(w, "error marshalling json", err_marshalling_json, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

// Only for use in dev environments to facilitate testing
func (cfg *apiConfig) ResetUserTable(w http.ResponseWriter, r *http.Request) {
	if cfg.environment != "dev" {
		custom_errors.ReturnErrorWithMessage(w, "incorrect environment", nil, 401)
	}
	err_resetting_db := cfg.database.DeleteAllUsers(r.Context())
	if err_resetting_db != nil {
		custom_errors.ReturnErrorWithMessage(w, "error resetting database", err_resetting_db, 500)
	}
}
