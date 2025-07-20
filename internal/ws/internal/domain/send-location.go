package domain

import (
	"fmt"
)

type SendLocation struct {
	deliveryDAO DeliveryDAO
}

func NewSendLocation(deliveryDAO DeliveryDAO) *SendLocation {
	return &SendLocation{
		deliveryDAO: deliveryDAO,
	}
}
func (r *SendLocation) Execute(inut *SendLocationInput) error {
	location := NewLocation(inut.Latitude, inut.Longitude)

	err := r.deliveryDAO.UpdateLocation(location, uint(inut.Delivery))
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
