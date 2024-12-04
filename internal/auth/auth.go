package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_password, err_hashing_password := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err_hashing_password != nil {
		return "", err_hashing_password
	}
	return string(hashed_password), nil
}

func CheckPassword(password, hashed_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	if len(tokenSecret) == 0 || userID == uuid.Nil {
		return "", fmt.Errorf("error with secret or user ID")
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "budget_buddy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   userID.String(),
	})
	tokenString, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func CheckJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err_parsing_claims := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err_parsing_claims != nil {
		return uuid.Nil, err_parsing_claims
	}
	idString, err_getting_id_from_token := token.Claims.GetSubject()
	if err_getting_id_from_token != nil {
		return uuid.Nil, err_getting_id_from_token
	}
	userId, err_parsing_id_str := uuid.Parse(idString)
	if err_parsing_id_str != nil {
		return uuid.Nil, err_parsing_id_str
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authToken := headers.Get("Authorization")
	if authToken == "" {
		return "", fmt.Errorf("no token provided")
	}
	token := strings.Split(authToken, " ")
	if len(token) != 2 {
		return "", fmt.Errorf("token error")
	}
	return token[1], nil
}
