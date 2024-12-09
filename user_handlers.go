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
	ID           uuid.UUID `json:"id"`
	Created_At   time.Time `json:"created_at"`
	Updated_At   time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Is_Admin     bool      `json:"is_admin"`
	First_Name   string    `json:"first_name"`
	Last_Name    string    `json:"last_name"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type AddUser struct {
	Email      string `json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Password   string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	new_user_info := AddUser{}
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
	w.WriteHeader(201)
}

func (cfg *apiConfig) UserLogin(w http.ResponseWriter, r *http.Request) {
	loginInfo := UserLogin{}
	decoder := json.NewDecoder(r.Body)
	err_decoding_body := decoder.Decode(&loginInfo)
	if err_decoding_body != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_decoding_body, 401)
		return
	}
	dbUser, err_retrieving_usr := cfg.database.GetUserByEmail(r.Context(), loginInfo.Email)
	if err_retrieving_usr != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_retrieving_usr, 401)
		return
	}
	err_checking_password := auth.CheckPassword(dbUser.HashedPassword, loginInfo.Password)
	if err_checking_password != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_checking_password, 500)
		return
	}
	newJWT, err_making_jwt := auth.MakeJWT(dbUser.ID, cfg.jwtSecret)
	if err_making_jwt != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_checking_password, 500)
		return
	}
	refreshToken, err_making_refresh := auth.MakeRefreshToken()
	if err_making_refresh != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 500)
	}
	cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	})
	loggedInUser := User{
		ID:           dbUser.ID,
		Created_At:   dbUser.CreatedAt,
		Updated_At:   dbUser.UpdatedAt,
		Email:        dbUser.Email,
		First_Name:   dbUser.FirstName,
		Last_Name:    dbUser.LastName,
		Token:        newJWT,
		RefreshToken: refreshToken,
	}
	jsonData, err_marshalling_json := json.Marshal(loggedInUser)
	if err_marshalling_json != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_marshalling_json, 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}
