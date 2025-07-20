package domain

import (
	"errors"
	"fmt"
)

type CreateDelivery struct {
	deliveryRepository DeliveryRepository
	userDao            UserDao
}

func NewCreateDelivery(deliveryRepository DeliveryRepository, userDao UserDao) *CreateDelivery {
	return &CreateDelivery{
		deliveryRepository: deliveryRepository,
		userDao:            userDao,
	}
}
func (r *CreateDelivery) Execute(inut *CreateDeliveryInput) error {

	driver, err := r.userDao.FindDriverById(inut.DriverId)
	if err != nil {
		return fmt.Errorf("error Create delivery: %w", err)
	}

	if driver == nil {
		return errors.New("error Create delivery: driver not found")
	}

	company, err := r.userDao.FindCompanyById(inut.CompanyId)
	if err != nil {
		return fmt.Errorf("error Create delivery: %w", err)
	}

	if company == nil {
		return errors.New("error Create delivery: company not found")
	}

	delivery := NewDelivery(inut.CompanyId, inut.DriverId, inut.Recipient, inut.Address)
	err = r.deliveryRepository.Save(delivery)
	if err != nil {
		return fmt.Errorf("error Create delivery: %w", err)
	}
	return nil

}

type CreateDeliveryInput struct {
	CompanyId uint
	DriverId  uint
	Recipient string
	Address   string
}
