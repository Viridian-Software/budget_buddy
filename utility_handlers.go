package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/Viridian-Software/budget_buddy/internal/auth"
	"github.com/Viridian-Software/budget_buddy/internal/constants"
	"github.com/google/uuid"
)

func (cfg *apiConfig) ServerRunningHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget buddy server is running"))
}

func (cfg *apiConfig) UserAuthentication(r *http.Request) (uuid.UUID, error) {
	jwt, err_invalid_bearer_token := auth.GetBearerToken(r.Header)
	if err_invalid_bearer_token != nil {
		return uuid.Nil, errors.New(constants.BEARER_TOKEN_ERROR)
	}
	userID, err_invalid_jwt := auth.CheckJWT(jwt, cfg.jwtSecret)
	if err_invalid_jwt != nil {
		return uuid.Nil, errors.New(constants.JWT_ERROR)
	}
	return userID, nil
}

func StringToFloat(value string) (float64, error) {
	value = RemoveNonNumeric(value)
	if dollarValue, err := strconv.ParseFloat(value, 64); err != nil {
		return 0.0, errors.New("unable to properly convert string")
	} else {
		return dollarValue, nil
	}
}

func RemoveNonNumeric(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsDigit(char) || char == '.' {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
