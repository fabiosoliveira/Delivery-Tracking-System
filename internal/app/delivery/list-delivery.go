package delivery

import (
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type ListDelivery struct {
	deliveryRepository domain.DeliveryRepository
	userRepository     domain.UserRepository
}

func NewListDelivery(deliveryRepository domain.DeliveryRepository, userRepository domain.UserRepository) *ListDelivery {
	return &ListDelivery{
		deliveryRepository: deliveryRepository,
		userRepository:     userRepository,
	}
}

func (r *ListDelivery) Execute(companyId int) ([]ListDeliveryOutput, error) {
	deliveries, err := r.deliveryRepository.ListDeliveryByCompanyId(companyId)
	if err != nil {
		return nil, fmt.Errorf("error listing deliveries: %w", err)
	}

	var deliveriesOutput []ListDeliveryOutput
	for _, delivery := range deliveries {
		driver, err := r.userRepository.FindById(delivery.Driver_id())
		if err != nil {
			return nil, fmt.Errorf("error listing deliveries: %w", err)
		}
		deliveriesOutput = append(deliveriesOutput, ListDeliveryOutput{delivery.Id(), *delivery.Status().String(), driver.Name(), delivery.Recipient(), delivery.Address()})
	}

	return deliveriesOutput, nil
}

type ListDeliveryOutput struct {
	Id        uint
	Status    string
	Driver    string
	Recipient string
	Address   string
}
