package driver

import (
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type ListDrivers struct {
	userRepository domain.UserRepository
}

func NewListDrivers(repository domain.UserRepository) *ListDrivers {
	return &ListDrivers{
		userRepository: repository,
	}
}

func (r *ListDrivers) Execute(userId int) ([]ListDriversOutput, error) {
	drivers, err := r.userRepository.ListDriversByCompanyId(userId)
	if err != nil {
		return nil, fmt.Errorf("error listing drivers: %w", err)
	}

	var driversOutput []ListDriversOutput
	for _, driver := range drivers {
		driversOutput = append(driversOutput, ListDriversOutput{int(driver.ID()), driver.Name(), driver.Email()})
	}

	return driversOutput, nil
}

type ListDriversOutput struct {
	ID    int
	Name  string
	Email string
}
