package entity

import "testing"

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber PhoneNumber
		expected    error
	}{
		{`Phone number must start with "+"`, "1234567890", ErrPhoneNumberPrefix},                 // Test case 1: Phone number must start with "+"
		{`Phone number length exceeds 15 characters`, "+1234567890123456", ErrPhoneNumberLength}, // Test case 2: Phone number length exceeds 15 characters
		{`Valid country code like "591"`, "+5911234567890", nil},                                 // Test case 3: Phone number with country code like "591"
		{`Valid country code like "1-246"`, "+1-246123456789", nil},
		{`Invalid country code`, "+0-242323232", ErrPhoneNumberCountryCode}, // Test case 3: Phone number with country code like "1-246"
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.phoneNumber.Valid()
			if result != test.expected {
				t.Errorf("For %s, expected %v, got %v", test.name, test.expected, result)
			}
		})
	}
}
