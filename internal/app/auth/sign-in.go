package auth

import (
	"errors"
	"fmt"
	"strconv"

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

func (su *SignIn) Execute(input *SignInInput) (*SignInOutput, error) {
	user, err := su.userRepository.FindByEmail(&input.Email)
	if err != nil {
		return nil, fmt.Errorf("error signing: %w", err)
	}

	if user == nil {
		return nil, errors.New("error signing: user not found")
	}

	if !user.VerifyPassword(input.Password) {
		return nil, errors.New("error signing: invalid password")
	}

	id := strconv.Itoa(int(user.ID()))
	return &SignInOutput{UserId: &id, UserType: user.UserType().String()}, nil
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	UserId   *string
	UserType *string
}
