package domain

import (
	"errors"
	"unicode/utf8"
)

type User interface {
	ID() uint
	Name() string
	Email() string
	Password() string
	VerifyPassword(password string) bool
	UserType() userType
}

type UserAbstract struct {
	id       uint
	name     string
	email    email
	password password
}

func newUser(name string, email string, password string) (User, error) {
	user := &UserAbstract{}

	errorName := user.SetName(name)
	errorEmail := user.SetEmail(email)
	errorPassword := user.SetPassword(password)

	err := ErrorsJoin(errorName, errorEmail, errorPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func restoreUser(id int, name string, _email string, passwordHash string) User {
	return &UserAbstract{
		id:       uint(id),
		name:     name,
		email:    email(_email),
		password: password(passwordHash),
	}
}

func (u *UserAbstract) SetName(name string) error {
	nameLength := utf8.RuneCountInString(name)

	if nameLength < 3 || nameLength > 30 {
		return errors.New("name must be between 3 and 30 characters")
	}

	u.name = name
	return nil
}

func (u *UserAbstract) SetEmail(e_mail string) error {
	_email := email(e_mail)

	err := _email.ValidateEmail()
	if err != nil {
		return err
	}

	u.email = _email
	return nil
}

func (u *UserAbstract) SetPassword(pass string) error {
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

func (u *UserAbstract) ID() uint {
	return u.id
}

func (u *UserAbstract) Name() string {
	return u.name
}

func (u *UserAbstract) Email() string {
	return string(u.email)
}

func (u *UserAbstract) Password() string {
	return string(u.password)
}

func (u *UserAbstract) VerifyPassword(password string) bool {
	return u.password.VerifyPassword(password)
}

func (u *UserAbstract) UserType() userType {
	panic("not implemented")
}
