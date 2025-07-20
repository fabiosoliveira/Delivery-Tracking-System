package domain

import (
	"errors"
	"fmt"
)

type SignUp struct {
	companyRepository CompanyRepository
}

func NewSignUp(repository CompanyRepository) *SignUp {
	return &SignUp{
		companyRepository: repository,
	}
}

func (su *SignUp) Execute(inut *SignUpInput) error {
	company, err := su.companyRepository.FindByEmail(&inut.Email)
	if err != nil {
		return fmt.Errorf("error signing up: %w", err)
	}

	if company != nil {
		return errors.New("error signing up: company already exists")
	}

	newCompany, err := NewCompany(inut.Name, inut.Email, inut.Password)
	if err != nil {
		return fmt.Errorf("error signing up: %w", err)
	}

	err = su.companyRepository.Save(newCompany)
	if err != nil {
		return fmt.Errorf("error signing up: %w", err)
	}
	return nil
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}
