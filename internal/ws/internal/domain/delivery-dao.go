package domain

type DeliveryDAO interface {
	UpdateLocation(location *Location, deliveryId uint) error
}

type Delivery struct {
	Id        uint
	Status    string
	Driver    string
	Recipient string
	Address   string
}
