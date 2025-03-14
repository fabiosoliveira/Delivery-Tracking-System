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

	newCompany, err := domain.NewCompany(inut.Name, inut.Email, inut.Password)
	if err != nil {
		return fmt.Errorf("error signing up: %w", err)
	}

	err = su.userRepository.Save(newCompany)
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
