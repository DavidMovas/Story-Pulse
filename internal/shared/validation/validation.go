package validation

import (
	"fmt"
	"net/mail"
	"story-pulse/contracts"
	"strings"

	"gopkg.in/validator.v2"
)

var (
	passwordMinLength       = 8
	emailMaxLength          = 127
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

func SetupValidators() {
	validators := []struct {
		name string
		fn   validator.ValidationFunc
	}{
		{"required", required},
		{"email", email},
		{"password", password},
		{"passwordOptional", passwordOptional},
		{"username", username},
		{"usernameOptional", usernameOptional},
		{"role", role},
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

func username(v interface{}, _ string) error {
	s, ok := v.(string)

	if !ok {
		return fmt.Errorf("username must be a string")
	}

	if len(s) < usernameMinLength {
		return fmt.Errorf("username must be at least %d characters long", usernameMinLength)
	}

	return nil
}

func usernameOptional(v interface{}, _ string) error {
	s, ok := v.(string)
	ps, psOk := v.(*string)

	if !ok && ps == nil {
		return nil
	}

	if psOk {
		s = *ps
		ok = true
	}

	if !ok {
		return fmt.Errorf("username must be a string")
	}

	if len(s) < usernameMinLength {
		return fmt.Errorf("username must be at least %d characters long", usernameMinLength)
	}

	return nil
}

func password(v interface{}, _ string) error {
	s, ok := v.(string)

	if !ok {
		return fmt.Errorf("password must be a string")
	}

	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
	}

	for _, entry := range passwordRequiredEntries {
		if !strings.ContainsAny(s, entry.chars) {
			return fmt.Errorf("password must contain at least one of the following required entries: %s", entry.name)
		}
	}

	return nil
}

func passwordOptional(v interface{}, _ string) error {
	s, ok := v.(string)
	ps, psOk := v.(*string)

	if !ok && ps == nil {
		return nil
	}

	if psOk {
		s = *ps
		ok = true
	}

	if !ok {
		return fmt.Errorf("password must be a string")
	}

	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
	}

	for _, entry := range passwordRequiredEntries {
		if !strings.ContainsAny(s, entry.chars) {
			return fmt.Errorf("password must contain at least one of the following required entries: %s", entry.name)
		}
	}

	return nil
}

func email(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	if len(s) > emailMaxLength {
		return fmt.Errorf("email must be at most %d characters long", emailMaxLength)
	}

	_, err := mail.ParseAddress(s)
	return err
}

func required(v interface{}, _ string) error {
	s, ok := v.(string)
	if ok && s == "" {
		return fmt.Errorf("must not be empty")
	}
	return nil
}

func role(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("role must be a string")
	}

	if !contracts.ValidateRole(contracts.Role(s)) {
		return fmt.Errorf("invalid role")
	}
	return nil
}
