package helpers

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) error {

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email")
	}

	parts := strings.Split(email, "@")

	// if there is no content before or after the @ symbol
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("invalid email")
	}

	afterAtSymbol := strings.Split(parts[1], ".")
	// if there is a dot after the @ symbol
	if len(afterAtSymbol) == 0 {
		return fmt.Errorf("invalid email")
	}

	// if there is content after the last dot(.)
	if len(afterAtSymbol[len(afterAtSymbol)-1]) == 0 {
		return fmt.Errorf("invalid email")
	}

	return nil
}

func ValidatePasswords(storedPassword, givenPassword []byte) error {

	if err := bcrypt.CompareHashAndPassword(storedPassword, givenPassword); err != nil {
		return err
	}
	return nil
}
