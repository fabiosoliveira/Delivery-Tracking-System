package domain

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type password string

func (p *password) ValidatePassword() error {
	if utf8.RuneCountInString(string(*p)) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(string(*p)) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(string(*p)) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(string(*p)) {
		return errors.New("password must contain at least one number")
	}
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(string(*p)) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// HashPassword generates a bcrypt hash for the given password.
func (p *password) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), 14)
	if err != nil {
		return err
	}

	*p = password(bytes)

	return nil

}

// VerifyPassword verifies if the given password matches the stored hash.
func (p *password) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(password))
	return err == nil
}
