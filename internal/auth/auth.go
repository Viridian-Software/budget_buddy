package auth

import "golang.org/x/crypto/bcrypt"

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
