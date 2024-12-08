package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Viridian-Software/budget_buddy/internal/custom_errors"
	"github.com/Viridian-Software/budget_buddy/internal/database"
	"github.com/google/uuid"
)

type Account struct {
	ID              uuid.UUID `json:"id"`
	Account_Name    string    `json:"account_name"`
	Current_Balance float64   `json:"current_balance"`
	User_ID         uuid.UUID `json:"user_id"`
	Account_Type    string    `json:"account_type"`
	Created_At      time.Time `json:"created_at"`
	Updated_at      time.Time `json:"updated_at"`
}

func (cfg *apiConfig) AddAccountHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	account_to_add := Account{}
	err_decoding_req := decoder.Decode(&account_to_add)
	if err_decoding_req != nil {
		custom_errors.ReturnErrorWithMessage(w, "error decoding json:", err_decoding_req, 500)
		return
	}
	dbAccount, err_adding_account := cfg.database.CreateNewAccount(r.Context(), database.CreateNewAccountParams{
		AccountName:    account_to_add.Account_Name,
		CurrentBalance: strconv.FormatFloat(account_to_add.Current_Balance, 'f', 10, 64),
		AccountType:    account_to_add.Account_Type,
		UserID:         account_to_add.User_ID,
	})
	if err_adding_account != nil {
		custom_errors.ReturnErrorWithMessage(w, "could not create new account", err_adding_account, 500)
		return
	}
	jsonData, err_marshalling_json := json.Marshal(dbAccount)
	if err_marshalling_json != nil {
		custom_errors.ReturnErrorWithMessage(w, "error creating json", err_marshalling_json, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (cfg *apiConfig) GetAllUserAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err_parsing_userID := uuid.Parse(r.URL.Query().Get("userID"))
	if err_parsing_userID != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 401)
		return
	}
	userAccount, err_retrieving_acc := cfg.database.GetAccountsByUser(r.Context(), userID)
	if err_retrieving_acc != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 404)
		return
	}
	accountSlice := []Account{}
	for _, value := range userAccount {
		num, err := strconv.ParseFloat(value.CurrentBalance, 64)
		if err != nil {
			custom_errors.ReturnErrorWithMessage(w, "", nil, 500)
			return
		}
		accountSlice = append(accountSlice, Account{
			ID:              value.ID,
			Account_Name:    value.AccountName,
			Current_Balance: num,
			User_ID:         value.UserID,
			Account_Type:    value.AccountType,
			Created_At:      value.CreatedAt,
			Updated_at:      value.UpdatedAt,
		})
	}
	jsonData, err_marshalling_json := json.Marshal(accountSlice)
	if err_marshalling_json != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 500)
	}
	w.WriteHeader(200)
	w.Write(jsonData)
}
