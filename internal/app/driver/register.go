package driver

import (
	"errors"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type Register struct {
	userRepository domain.UserRepository
}

func NewRegister(repository domain.UserRepository) *Register {
	return &Register{
		userRepository: repository,
	}
}

func (r *Register) Execute(inut *RegisterInput) error {
	driver, err := r.userRepository.FindByEmail(&inut.Email)
	if err != nil {
		return fmt.Errorf("error register driver: %w", err)
	}

	if driver != nil {
		return errors.New("error register driver: driver already exists")
	}

	driver, err = domain.NewDriver(inut.Name, inut.Email, inut.Password, inut.CompanyId)
	if err != nil {
		return fmt.Errorf("error register driver: %w", err)
	}

	err = r.userRepository.Save(driver)
	if err != nil {
		return fmt.Errorf("error register driver: %w", err)
	}
	return nil
}

type RegisterInput struct {
	Name      string
	Email     string
	Password  string
	CompanyId uint
}
