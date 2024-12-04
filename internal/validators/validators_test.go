package validators

import "testing"

func TestPasswordMeetsRequirements(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Password below minimum length",
			password: "short",
			expected: false,
		},
		{
			name:     "Password exceeds maximum length",
			password: "thispasswordistoolong",
			expected: false,
		},
		{
			name:     "Password at minimum length",
			password: "minimum",
			expected: false,
		},
		{
			name:     "Password at maximum length",
			password: "fifteenletter",
			expected: true,
		},
		{
			name:     "Password within valid length",
			password: "valid123",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidatePassword(test.password)
			if result != test.expected {
				t.Errorf("PasswordMeetsRequirements(%q) = %v; want %v", test.password, result, test.expected)
			}
		})
	}
}
