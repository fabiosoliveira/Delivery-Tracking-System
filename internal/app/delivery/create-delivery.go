package delivery

import (
	"errors"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type CreateDelivery struct {
	deliveryRepository domain.DeliveryRepository
	userRepository     domain.UserRepository
}

func NewCreateDelivery(deliveryRepository domain.DeliveryRepository, userRepository domain.UserRepository) *CreateDelivery {
	return &CreateDelivery{
		deliveryRepository: deliveryRepository,
		userRepository:     userRepository,
	}
}
func (r *CreateDelivery) Execute(inut *CreateDeliveryInput) error {

	driver, err := r.userRepository.FindById(inut.DriverId)
	if err != nil {
		return fmt.Errorf("error Create delivery: %w", err)
	}

	if driver == nil {
		return errors.New("error Create delivery: driver not found")
	}

	company, err := r.userRepository.FindById(inut.CompanyId)
	if err != nil {
		return fmt.Errorf("error Create delivery: %w", err)
	}

	if company == nil {
		return errors.New("error Create delivery: company not found")
	}

	delivery := domain.NewDelivery(inut.CompanyId, inut.DriverId, inut.Recipient, inut.Address)
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
