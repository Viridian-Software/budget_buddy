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

type Transaction struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	User_ID    uuid.UUID `json:"user_id"`
	Account_ID uuid.UUID `json:"account_id"`
	Amount     float64   `json:"amount"`
}

type TransactionCreated struct {
	Transaction_Details Transaction `json:"transaction_details"`
	Account_Details     Account     `json:"account_details"`
}

func (cfg *apiConfig) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	user_id, err_authenticating_user := cfg.UserAuthentication(r)
	if err_authenticating_user != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_authenticating_user, http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	new_transaction := Transaction{}
	err_decoding_request := decoder.Decode(&new_transaction)
	if err_decoding_request != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_decoding_request, http.StatusInternalServerError)
		return
	}
	if user_id != new_transaction.User_ID {
		custom_errors.ReturnErrorWithMessage(w, "", nil, http.StatusUnauthorized)
		return
	}
	user_account, err_fetching_account := cfg.database.GetAccount(r.Context(), new_transaction.Account_ID)
	if err_fetching_account != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_fetching_account, http.StatusUnauthorized)
		return
	}
	current_balance, err_parsing_float := strconv.ParseFloat(user_account.CurrentBalance, 64)
	if err_parsing_float != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_parsing_float, http.StatusInternalServerError)
		return
	}
	new_balance := current_balance + new_transaction.Amount
	updated_account, err_updating_account := cfg.database.UpdateBalance(r.Context(), database.UpdateBalanceParams{
		ID:             new_transaction.Account_ID,
		CurrentBalance: strconv.FormatFloat(new_balance, 'f', 10, 64),
	})
	if err_updating_account != nil {
		custom_errors.ReturnErrorWithMessage(w, "error: unable to update account", err_updating_account, http.StatusInternalServerError)
		return
	}
	return_information := TransactionCreated{}
	return_information.Account_Details = Account{
		ID:              updated_account.ID,
		Account_Name:    updated_account.AccountName,
		Current_Balance: new_balance,
		User_ID:         updated_account.UserID,
		Account_Type:    updated_account.AccountType,
		Created_At:      updated_account.CreatedAt,
		Updated_at:      updated_account.UpdatedAt,
	}

	db_transaction, err_adding_transaction := cfg.database.CreateTransaction(r.Context(), database.CreateTransactionParams{
		ID:        new_transaction.ID,
		UserID:    user_id,
		AccountID: new_transaction.Account_ID,
		Amount:    strconv.FormatFloat(new_transaction.Amount, 'f', 10, 64),
	})
	if err_adding_transaction != nil {
		custom_errors.ReturnErrorWithMessage(w, "error: unable to add transaction", err_adding_transaction, http.StatusInternalServerError)
		return
	}
	return_information.Transaction_Details = Transaction{
		ID:         db_transaction.ID,
		Created_At: db_transaction.CreatedAt,
		User_ID:    db_transaction.UserID,
		Account_ID: db_transaction.AccountID,
		Amount:     new_transaction.Amount,
	}
	json_transaction, err_marshalling_transaction := json.Marshal(return_information)
	if err_marshalling_transaction != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_transaction)
	w.WriteHeader(http.StatusOK)
}
