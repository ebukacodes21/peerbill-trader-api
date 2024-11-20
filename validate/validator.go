package validate

import "fmt"

func ValidateString(value string, min int, max int) error {
	n := len(value)

	if n < min || n > max {
		return fmt.Errorf("must contain from %d-%d characters", min, max)
	}
	return nil
}

func ValidateFirstname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	return nil
}

func ValidateLastname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	return nil
}
