package domain

import (
	"errors"
	"fmt"
	"strconv"
)

type SignIn struct {
	companyRepository CompanyRepository
}

func NewSignIn(repository CompanyRepository) *SignIn {
	return &SignIn{
		companyRepository: repository,
	}
}

func (si *SignIn) Execute(input *SignInInput) (*SignInOutput, error) {
	company, err := si.companyRepository.FindByEmail(&input.Email)
	if err != nil {
		return nil, fmt.Errorf("error signing: %w", err)
	}

	if company == nil {
		return nil, errors.New("error signing: company not found")
	}

	if !company.VerifyPassword(input.Password) {
		return nil, errors.New("error signing: invalid password")
	}

	id := strconv.Itoa(int(company.ID()))
	return &SignInOutput{UserId: id, UserType: "company"}, nil
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	UserId   string
	UserType string
}
