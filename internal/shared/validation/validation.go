package validation

import (
	"fmt"
	"net/mail"
	"strings"
)

type Type string

const (
	Email    = "email"
	Username = "username"
	Password = "password"
)

var (
	emailMaxLength          = 127
	passwordMinLength       = 8
	usernameMinLength       = 3
	passwordSpecialChars    = "!#$%&'*+/=?^_`{|}~@"
	passwordRequiredEntries = []struct {
		name  string
		chars string
	}{
		{"lowercase", "abcdefghijklmnopqrstuvwxyz"},
		{"uppercase", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"numbers", "0123456789"},
		{"special character (" + passwordSpecialChars + ")", passwordSpecialChars},
	}
)

func Validate(validType string, value any, required bool) error {
	s, ok := value.(string)
	if !ok && required {
		return fmt.Errorf("%s is required", validType)
	}

	if !ok {
		return fmt.Errorf("%s is %s", validType, value)
	}

	switch validType {
	case Email:
		return validateEmail(s, required)
	case Username:
		return validateUsername(s, required)
	case Password:
		return validatePassword(s, required)
	default:
		return nil
	}
}

func validateEmail(email string, required bool) error {
	if required && email == "" {
		return fmt.Errorf("email is required")
	} else if !required && email == "" {
		return nil
	}

	if len(email) > emailMaxLength {
		return fmt.Errorf("email is too long")
	}

	_, err := mail.ParseAddress(email)
	return err
}

func validateUsername(username string, required bool) error {
	if required && username == "" {
		return fmt.Errorf("username is required")
	} else if !required && username == "" {
		return nil
	}

	if len(username) < usernameMinLength {
		return fmt.Errorf("username too short")
	}

	return nil
}

func validatePassword(password string, required bool) error {
	if required && password == "" {
		return fmt.Errorf("password is required")
	} else if !required && password == "" {
		return nil
	}

	if len(password) < passwordMinLength {
		return fmt.Errorf("password too short")
	}

	for _, entry := range passwordRequiredEntries {
		if !strings.ContainsAny(password, entry.chars) {
			return fmt.Errorf("password must contain at least one of the following required entries: %s", entry.name)
		}
	}

	return nil
}
