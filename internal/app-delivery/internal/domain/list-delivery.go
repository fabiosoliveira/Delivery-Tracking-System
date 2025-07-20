package domain

import (
	"fmt"
)

type ListDelivery struct {
	deliveryRepository DeliveryRepository
	driverRepository   DriverRepository
}

func NewListDelivery(deliveryRepository DeliveryRepository, driverRepository DriverRepository) *ListDelivery {
	return &ListDelivery{
		deliveryRepository: deliveryRepository,
		driverRepository:   driverRepository,
	}
}

func (r *ListDelivery) Execute(companyId int) ([]ListDeliveryOutput, error) {
	deliveries, err := r.deliveryRepository.ListDeliveryByCompanyId(companyId)
	if err != nil {
		return nil, fmt.Errorf("error listing deliveries: %w", err)
	}

	var deliveriesOutput []ListDeliveryOutput
	for _, delivery := range deliveries {
		driver, err := r.driverRepository.FindById(delivery.DriverId)
		if err != nil {
			return nil, fmt.Errorf("error listing deliveries: %w", err)
		}
		deliveriesOutput = append(deliveriesOutput, ListDeliveryOutput{delivery.Id, *delivery.Status.String(), driver.Name(), delivery.Recipient, delivery.Address})
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
