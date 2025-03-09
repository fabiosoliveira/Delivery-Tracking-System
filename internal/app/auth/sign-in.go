package auth

import (
	"errors"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type SignIn struct {
	userRepository domain.UserRepository
}

func NewSignIn(repository domain.UserRepository) *SignIn {
	return &SignIn{
		userRepository: repository,
	}
}

func (su *SignIn) Execute(inut *SignInInput) error {
	user, err := su.userRepository.FindByEmail(&inut.Email)
	if err != nil {
		return fmt.Errorf("error signing: %w", err)
	}

	if user == nil {
		return errors.New("error signing: user not found")
	}

	// TODO: implementar autenticação

	return nil
}

type SignInInput struct {
	Name     string
	Email    string
	Password string
}
