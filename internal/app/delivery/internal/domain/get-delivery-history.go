package domain

type GetDeliveryHistory struct {
	deliveryRepository DeliveryRepository
}

func NewGetDeliveryHistory(deliveryRepository DeliveryRepository) *GetDeliveryHistory {
	return &GetDeliveryHistory{
		deliveryRepository: deliveryRepository,
	}
}

func (g *GetDeliveryHistory) Execute(deliveryID int) ([]Output, error) {
	locations, err := g.deliveryRepository.FindLocationsByDeliveryID(deliveryID)
	if err != nil {
		return nil, err
	}

	var output []Output
	for _, location := range locations {
		output = append(output, Output{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		})
	}
	return output, nil
}

type Output struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
