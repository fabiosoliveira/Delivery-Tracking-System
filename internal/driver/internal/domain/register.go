package domain

import (
	"errors"
	"fmt"
)

type Register struct {
	driverRepository DriverRepository
}

func NewRegister(repository DriverRepository) *Register {
	return &Register{
		driverRepository: repository,
	}
}

func (r *Register) Execute(inut *RegisterInput) error {
	driver, err := r.driverRepository.FindByEmail(&inut.Email)
	if err != nil {
		return fmt.Errorf("error register driver: %w", err)
	}

	if driver != nil {
		return errors.New("error register driver: driver already exists")
	}

	driver, err = NewDriver(inut.Name, inut.Email, inut.Password, inut.CompanyId)
	if err != nil {
		return fmt.Errorf("error register driver: %w", err)
	}

	err = r.driverRepository.Save(driver)
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
