package delivery

import (
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type SendLocation struct {
	deliveryRepository domain.DeliveryRepository
}

func NewSendLocation(deliveryRepository domain.DeliveryRepository) *SendLocation {
	return &SendLocation{
		deliveryRepository: deliveryRepository,
	}
}
func (r *SendLocation) Execute(inut *SendLocationInput) error {
	location := domain.NewLocation(inut.Latitude, inut.Longitude)

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
