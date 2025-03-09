package auth

import (
	"errors"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type SignUp struct {
	userRepository domain.UserRepository
}

func NewSignUp(repository domain.UserRepository) *SignUp {
	return &SignUp{
		userRepository: repository,
	}
}

func (su *SignUp) Execute(inut *SignUpInput) error {
	company, err := su.userRepository.FindByEmail(&inut.Email)
	if err != nil {
		return fmt.Errorf("error signing up: %w", err)
	}

	if company != nil {
		return errors.New("error signing up: company already exists")
	}

	company = &domain.User{
		Name:     inut.Name,
		Email:    inut.Email,
		Password: inut.Password,
		UserType: domain.UserTypeCompany,
	}

	err = su.userRepository.Save(company)
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
