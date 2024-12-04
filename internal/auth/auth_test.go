package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successfully hashes a password", func(t *testing.T) {
		password := "securepassword123"
		hashedPassword, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Ensure hashed password is not empty and doesn't match the raw password
		if hashedPassword == "" {
			t.Fatalf("expected hashed password to not be empty")
		}
		if hashedPassword == password {
			t.Fatalf("hashed password should not be the same as the raw password")
		}
	})

	t.Run("Handles empty password input", func(t *testing.T) {
		password := ""
		hashedPassword, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error for empty password, got %v", err)
		}

		if hashedPassword == "" {
			t.Fatalf("expected hashed password to not be empty for empty input")
		}
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("Successfully validates a correct password", func(t *testing.T) {
		password := "securepassword123"
		hashedPassword, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error while hashing password, got %v", err)
		}

		err = CheckPassword(password, hashedPassword)
		if err != nil {
			t.Fatalf("expected password to validate, got error: %v", err)
		}
	})

	t.Run("Fails validation for an incorrect password", func(t *testing.T) {
		password := "securepassword123"
		incorrectPassword := "wrongpassword"
		hashedPassword, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error while hashing password, got %v", err)
		}

		err = CheckPassword(incorrectPassword, hashedPassword)
		if err == nil {
			t.Fatalf("expected validation to fail for incorrect password, but it succeeded")
		}
	})

	t.Run("Fails validation for a malformed hash", func(t *testing.T) {
		password := "securepassword123"
		malformedHash := "notarealhash"

		err := CheckPassword(password, malformedHash)
		if err == nil {
			t.Fatalf("expected validation to fail for a malformed hash, but it succeeded")
		}
	})
}
