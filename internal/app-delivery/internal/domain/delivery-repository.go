package domain

type DeliveryRepository interface {
	// Save(delivery *Delivery) error
	ListDeliveryByCompanyId(id int) ([]Delivery, error)
	// ListDeliveryByDriverId(driverId int) ([]Delivery, error)
	// FindById(deliveryId uint) (*Delivery, error)
	// UpdateLocation(location *Location, deliveryId uint) error
	// FindLocationsByDeliveryID(deliveryID int) ([]Location, error)
}
