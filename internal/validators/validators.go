package validators

import (
	"regexp"
	"strings"
)

func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 15 {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	// Trim leading/trailing spaces.
	email = strings.TrimSpace(email)

	// Ensure there is exactly one '@'.
	if strings.Count(email, "@") != 1 {
		return false
	}

	// Ensure there are no consecutive dots in the email.
	if strings.Contains(email, "..") {
		return false
	}

	// Ensure email matches the regular expression.
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
