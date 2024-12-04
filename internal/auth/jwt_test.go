package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMakeJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	token, err := MakeJWT(userID, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token to verify the claims
	parsedToken, _ := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	assert.True(t, ok)
	assert.Equal(t, "budget_buddy", claims.Issuer)
	assert.Equal(t, userID.String(), claims.Subject)
	assert.WithinDuration(t, time.Now().UTC().Add(time.Hour), claims.ExpiresAt.Time, time.Second*5)
}

func TestMakeJWT_InvalidInputs(t *testing.T) {
	secret := "test_secret"

	// Test with empty secret
	token, err := MakeJWT(uuid.New(), "")
	assert.Error(t, err)
	assert.Equal(t, "", token)

	// Test with nil user ID
	token, err = MakeJWT(uuid.Nil, secret)
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestCheckJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	// Create a valid JWT
	token, _ := MakeJWT(userID, secret)

	// Check the JWT
	parsedUserID, err := CheckJWT(token, secret)

	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestCheckJWT_InvalidToken(t *testing.T) {
	secret := "test_secret"
	invalidToken := "this.is.not.a.valid.token"

	userID, err := CheckJWT(invalidToken, secret)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, userID)
}

func TestCheckJWT_InvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"
	wrongSecret := "wrong_secret"

	// Create a valid JWT with the correct secret
	token, _ := MakeJWT(userID, secret)

	// Try to validate the JWT with the wrong secret
	parsedUserID, err := CheckJWT(token, wrongSecret)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedUserID)
}

func TestCheckJWT_MissingSubject(t *testing.T) {
	secret := "test_secret"

	// Create a JWT with missing subject
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "budget_buddy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	// Check the JWT
	userID, err := CheckJWT(tokenString, secret)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, userID)
}

func TestCheckJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "budget_buddy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour)),
		Subject:   userID.String(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	// Check the expired JWT
	parsedUserID, err := CheckJWT(tokenString, secret)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedUserID)
}
