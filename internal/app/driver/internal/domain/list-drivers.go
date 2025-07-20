package domain

import (
	"fmt"
)

type ListDrivers struct {
	driverRepository DriverRepository
}

func NewListDrivers(repository DriverRepository) *ListDrivers {
	return &ListDrivers{
		driverRepository: repository,
	}
}

func (r *ListDrivers) Execute(userId int) ([]ListDriversOutput, error) {
	drivers, err := r.driverRepository.ListDriversByCompanyId(userId)
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
