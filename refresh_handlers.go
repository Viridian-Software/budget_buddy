package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Viridian-Software/budget_buddy/internal/auth"
	"github.com/Viridian-Software/budget_buddy/internal/custom_errors"
)

func (cfg *apiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err_fetching_token := auth.GetBearerToken(r.Header)
	if err_fetching_token != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_fetching_token, 401)
		return
	}
	dbToken, err_retrieving_token_from_db := cfg.database.CheckRefreshToken(r.Context(), refreshToken)
	if err_retrieving_token_from_db != nil {
		custom_errors.ReturnErrorWithMessage(w, "", err_retrieving_token_from_db, 403)
		return
	}
	if time.Now().After(dbToken.ExpiresAt) || dbToken.RevokedAt.Valid {
		w.WriteHeader(401)
		return
	}
	newJWT, err_creating_jwt := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret)
	if err_creating_jwt != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 500)
		return
	}
	jsonData, err_marshalling_json := json.Marshal(struct {
		Token string
	}{Token: newJWT})
	if err_marshalling_json != nil {
		custom_errors.ReturnErrorWithMessage(w, "", nil, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (cfg *apiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	current_refresh_token, err_fetching_token_in_header := auth.GetBearerToken(r.Header)
	if err_fetching_token_in_header != nil {
		custom_errors.ReturnErrorWithMessage(w, "token error", nil, 403)
		return
	}
	cfg.database.RevokeToken(r.Context(), current_refresh_token)
	w.WriteHeader(204)
}
