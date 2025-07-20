package valueObject

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Password string

func (p *Password) ValidatePassword() error {
	var errs []error
	if utf8.RuneCountInString(string(*p)) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(string(*p)) {
		errs = append(errs, errors.New("password must contain at least one lowercase letter"))
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(string(*p)) {
		errs = append(errs, errors.New("password must contain at least one uppercase letter"))
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(string(*p)) {
		errs = append(errs, errors.New("password must contain at least one number"))
	}
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(string(*p)) {
		errs = append(errs, errors.New("password must contain at least one special character"))
	}

	return utils.ErrorsJoin(errs...)
}

// HashPassword generates a bcrypt hash for the given password.
func (p *Password) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), 14)
	if err != nil {
		return err
	}

	*p = Password(bytes)

	return nil

}

// VerifyPassword verifies if the given password matches the stored hash.
func (p *Password) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(password))
	return err == nil
}
