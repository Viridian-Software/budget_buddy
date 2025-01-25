package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	// Authenticate user
	userID, err := cfg.UserAuthentication(r)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "authentication failed", err, http.StatusUnauthorized)
		return
	}

	// Decode request body into a new transaction
	var newTransaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		custom_errors.ReturnErrorWithMessage(w, "invalid request body", err, http.StatusBadRequest)
		return
	}

	// Validate transaction user ID
	if userID != newTransaction.User_ID {
		custom_errors.ReturnErrorWithMessage(w, "unauthorized user", err, http.StatusUnauthorized)
		return
	}

	// Fetch user account details
	userAccount, err := cfg.database.GetAccount(r.Context(), newTransaction.Account_ID)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "account not found", err, http.StatusUnauthorized)
		return
	}

	// Parse current balance
	currentBalance, err := strconv.ParseFloat(userAccount.CurrentBalance, 64)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "invalid account balance", err, http.StatusInternalServerError)
		return
	}

	// Update account balance
	newBalance := currentBalance + newTransaction.Amount
	updatedAccount, err := cfg.database.UpdateBalance(r.Context(), database.UpdateBalanceParams{
		ID:             newTransaction.Account_ID,
		CurrentBalance: strconv.FormatFloat(newBalance, 'f', 10, 64),
	})
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "unable to update account balance", err, http.StatusInternalServerError)
		return
	}

	// Create a transaction record in the database
	dbTransaction, err := cfg.database.CreateTransaction(r.Context(), database.CreateTransactionParams{
		ID:        newTransaction.ID,
		UserID:    userID,
		AccountID: newTransaction.Account_ID,
		Amount:    strconv.FormatFloat(newTransaction.Amount, 'f', 10, 64),
	})
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "unable to create transaction", err, http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := TransactionCreated{
		Account_Details: Account{
			ID:              updatedAccount.ID,
			Account_Name:    updatedAccount.AccountName,
			Current_Balance: newBalance,
			User_ID:         updatedAccount.UserID,
			Account_Type:    updatedAccount.AccountType,
			Created_At:      updatedAccount.CreatedAt,
			Updated_at:      updatedAccount.UpdatedAt,
		},
		Transaction_Details: Transaction{
			ID:         dbTransaction.ID,
			Created_At: dbTransaction.CreatedAt,
			User_ID:    dbTransaction.UserID,
			Account_ID: dbTransaction.AccountID,
			Amount:     newTransaction.Amount,
		},
	}

	// Marshal response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Send status code before writing body
	if err := json.NewEncoder(w).Encode(response); err != nil {
		custom_errors.ReturnErrorWithMessage(w, "unable to encode response", err, http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	user_ID, err := cfg.UserAuthentication(r)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "authentication error", err, http.StatusUnauthorized)
		return
	}
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		custom_errors.ReturnErrorWithMessage(w, "error processing request", err, http.StatusInternalServerError)
		return
	}
	if user_ID != transaction.User_ID {
		custom_errors.ReturnErrorWithMessage(w, "unauthorized access", nil, http.StatusUnauthorized)
		return
	}
	trans, err := cfg.database.GetTransactionByID(r.Context(), transaction.ID)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "invalid transaction", err, 404)
		return
	}
	acc, err := cfg.database.GetAccount(r.Context(), trans.AccountID)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "invalid account", err, 404)
		return
	}
	currentAccountBalance, err1 := StringToFloat(acc.CurrentBalance)
	transactionAmount, err2 := StringToFloat(trans.Amount)
	if err1 != nil || err2 != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err, 500)
		return
	}
	newBalance := currentAccountBalance - transactionAmount
	_, err = cfg.database.UpdateBalance(r.Context(), database.UpdateBalanceParams{
		ID:             trans.AccountID,
		CurrentBalance: strconv.FormatFloat(newBalance, 'f', 10, 64),
	})
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "unable to update account balance", err, http.StatusInternalServerError)
		return
	}
	if err := cfg.database.DeleteTransaction(r.Context(), transaction.ID); err != nil {
		custom_errors.ReturnErrorWithMessage(w, "error processing request", nil, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) GetTransactionsForAccount(w http.ResponseWriter, r *http.Request) {
	_, err := cfg.UserAuthentication(r)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "authentication error", err, http.StatusUnauthorized)
		return
	}
	account_id, err_parsing_userID := uuid.Parse(r.PathValue("accountID"))
	if err_parsing_userID != nil {
		custom_errors.ReturnErrorWithMessage(w, "unauthorized", err_parsing_userID, http.StatusUnauthorized)
		return
	}
	transactions, err := cfg.database.GetAllTransactions(r.Context(), account_id)
	if err != nil {
		custom_errors.ReturnErrorWithMessage(w, "error retrieving account information", err, http.StatusUnauthorized)
		return
	}
	userTransactions := []Transaction{}
	for _, values := range transactions {
		values.Amount = strings.Trim(values.Amount, "$")
		amount, err := strconv.ParseFloat(values.Amount, 64)
		if err != nil {
			custom_errors.ReturnErrorWithMessage(w, "error fetching transaction amount", err, http.StatusInternalServerError)
			continue
		}
		userTransactions = append(userTransactions, Transaction{
			ID:         values.ID,
			Created_At: values.CreatedAt,
			User_ID:    values.UserID,
			Account_ID: values.AccountID,
			Amount:     amount,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userTransactions); err != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, http.StatusInternalServerError)
	}

}
