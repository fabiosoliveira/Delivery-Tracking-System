package domain

import (
	"errors"
	"unicode/utf8"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/internal/utils"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/internal/valueObject"
)

type Driver struct {
	id        uint
	name      string
	email     valueObject.Email
	password  valueObject.Password
	companyId uint
}

func NewDriver(name string, email string, password string, _companyId uint) (*Driver, error) {
	driver := &Driver{}

	errorName := driver.SetName(name)
	errorEmail := driver.SetEmail(email)
	errorPassword := driver.SetPassword(password)

	driver.companyId = _companyId

	err := utils.ErrorsJoin(errorName, errorEmail, errorPassword)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func RestoreDriver(id int, name string, _email string, passwordHash string, _companyId int) *Driver {
	return &Driver{
		id:        uint(id),
		name:      name,
		email:     valueObject.Email(_email),
		password:  valueObject.Password(passwordHash),
		companyId: uint(_companyId),
	}
}

func (u *Driver) SetName(name string) error {
	nameLength := utf8.RuneCountInString(name)

	if nameLength < 3 || nameLength > 30 {
		return errors.New("name must be between 3 and 30 characters")
	}

	u.name = name
	return nil
}

func (u *Driver) SetEmail(e_mail string) error {
	_email := valueObject.Email(e_mail)

	err := _email.ValidateEmail()
	if err != nil {
		return err
	}

	u.email = _email
	return nil
}

func (u *Driver) SetPassword(pass string) error {
	_password := valueObject.Password(pass)

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

func (u *Driver) ID() uint {
	return u.id
}

func (u *Driver) Name() string {
	return u.name
}

func (u *Driver) Email() string {
	return string(u.email)
}

func (u *Driver) Password() string {
	return string(u.password)
}

func (u *Driver) CompanyId() uint {
	return u.companyId
}

func (u *Driver) VerifyPassword(password string) bool {
	return u.password.VerifyPassword(password)
}
