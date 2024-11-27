package validate

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isUsernameValid = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, min int, max int) error {
	n := len(value)

	if n < min || n > max {
		return fmt.Errorf("must contain from %d-%d characters", min, max)
	}
	return nil
}

func ValidateFirstname(value string) error {
	return ValidateString(value, 3, 100)
}

func ValidateLastname(value string) error {
	return ValidateString(value, 3, 100)
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isUsernameValid(value) {
		return fmt.Errorf("must contain only letters, digits or underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 8, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}

func ValidatePhone(value string) error {
	return ValidateString(value, 11, 11)
}

func ValidateCountry(value string) error {
	return ValidateString(value, 1, 30)
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("value must be a positive integer")
	}
	return nil
}

func ValidateCode(value string) error {
	return ValidateString(value, 32, 128)
}

func ValidateFiat(value string) error {
	return ValidateString(value, 3, 4)
}

func ValidateCrypto(value string) error {
	return ValidateString(value, 3, 4)
}
