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

func (su *SignIn) Execute(inut *SignInInput) (*SignInOutput, error) {
	user, err := su.userRepository.FindByEmail(&inut.Email)
	if err != nil {
		return nil, fmt.Errorf("error signing: %w", err)
	}

	if user == nil {
		return nil, errors.New("error signing: user not found")
	}

	if user.Password != inut.Password {
		return nil, errors.New("error signing: invalid password")
	}

	return &SignInOutput{UserId: user.ID, UserType: user.UserType.String()}, nil
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	UserId   uint
	UserType *string
}
