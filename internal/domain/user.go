package domain

import (
	"errors"
	"unicode/utf8"
)

type User struct {
	id       uint
	name     string
	email    email
	password password
	userType userType
}

func NewUser(name string, email string, password string, userType userType) (*User, error) {
	user := &User{}

	errorName := user.SetName(name)
	errorEmail := user.SetEmail(email)
	errorPassword := user.SetPassword(password)
	user.userType = userType

	err := errors.Join(errorName, errorEmail, errorPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func RestoreUser(id int, name string, _email string, passwordHash string, _userType int) *User {
	return &User{
		id:       uint(id),
		name:     name,
		email:    email(_email),
		password: password(passwordHash),
		userType: userType(_userType),
	}
}

func (u *User) SetName(name string) error {
	nameLength := utf8.RuneCountInString(name)

	if nameLength < 3 || nameLength > 30 {
		return errors.New("name must be between 3 and 30 characters")
	}

	u.name = name
	return nil
}

func (u *User) SetEmail(e_mail string) error {
	_email := email(e_mail)

	err := _email.ValidateEmail()
	if err != nil {
		return err
	}

	u.email = _email
	return nil
}

func (u *User) SetPassword(pass string) error {
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

// func (u *User) SetUserType(userType userType) {
// 	u.userType = userType
// }

func (u *User) ID() uint {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return string(u.email)
}

func (u *User) Password() string {
	return string(u.password)
}

func (u *User) VerifyPassword(password string) bool {
	return u.password.VerifyPassword(password)
}

func (u *User) UserType() userType {
	return u.userType
}
