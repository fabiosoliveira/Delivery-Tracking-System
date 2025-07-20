package domain

import (
	"fmt"
)

type SendLocation struct {
	deliveryRepository DeliveryRepository
}

func NewSendLocation(deliveryRepository DeliveryRepository) *SendLocation {
	return &SendLocation{
		deliveryRepository: deliveryRepository,
	}
}
func (r *SendLocation) Execute(inut *SendLocationInput) error {
	location := NewLocation(inut.Latitude, inut.Longitude)

	err := r.deliveryRepository.UpdateLocation(location, uint(inut.Delivery))
	if err != nil {
		return fmt.Errorf("error sending location: %w", err)
	}
	return nil
}

type SendLocationInput struct {
	Delivery  int64
	Latitude  float64
	Longitude float64
}
