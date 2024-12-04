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
		custom_errors.HandleServerError(w, "error decoding json", err_decoding_response)
		return
	}
	if !validators.ValidatePassword(new_user_info.Password) {
		custom_errors.HandleServerError(w, "password does not meet requirements", nil)
	}
	hashed_password, err_hashing_password := auth.HashPassword(new_user_info.Password)
	if err_hashing_password != nil {
		custom_errors.HandleServerError(w, "error hashing password", err_hashing_password)
		return
	}
	dbUser, err_adding_user := cfg.database.AddUser(r.Context(), database.AddUserParams{
		Email:          new_user_info.Email,
		FirstName:      new_user_info.First_Name,
		LastName:       new_user_info.Last_Name,
		HashedPassword: hashed_password,
	})
	if err_adding_user != nil {
		custom_errors.HandleServerError(w, "error adding user", err_adding_user)
		return
	}
	jsonData, err_marshalling_json := json.Marshal(dbUser)
	if err_marshalling_json != nil {
		custom_errors.HandleServerError(w, "error marshalling json", err_marshalling_json)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}
