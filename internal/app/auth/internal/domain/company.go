package domain

import (
	"errors"
	"unicode/utf8"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/internal/utils"
)

type Company struct {
	id       uint
	name     string
	email    email
	password password
}

func NewCompany(name string, email string, password string) (*Company, error) {
	company := &Company{}

	errorName := company.SetName(name)
	errorEmail := company.SetEmail(email)
	errorPassword := company.SetPassword(password)

	err := utils.ErrorsJoin(errorName, errorEmail, errorPassword)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func RestoreCompany(id int, name string, _email string, passwordHash string) *Company {
	return &Company{
		id:       uint(id),
		name:     name,
		email:    email(_email),
		password: password(passwordHash),
	}
}

func (u *Company) SetName(name string) error {
	nameLength := utf8.RuneCountInString(name)

	if nameLength < 3 || nameLength > 30 {
		return errors.New("name must be between 3 and 30 characters")
	}

	u.name = name
	return nil
}

func (u *Company) SetEmail(e_mail string) error {
	_email := email(e_mail)

	err := _email.ValidateEmail()
	if err != nil {
		return err
	}

	u.email = _email
	return nil
}

func (u *Company) SetPassword(pass string) error {
	_password := password(pass)

	err := _password.ValidatePassword()
	if err != nil {
		return err
	}

	err = _password.HashPassword()
	if err != nil {
		return err
	}

	u.password = _password
	return nil
}

func (u *Company) ID() uint {
	return u.id
}

func (u *Company) Name() string {
	return u.name
}

func (u *Company) Email() string {
	return string(u.email)
}

func (u *Company) Password() string {
	return string(u.password)
}

func (u *Company) VerifyPassword(password string) bool {
	return u.password.VerifyPassword(password)
}
