package domain

type Delivery struct {
	Id        uint
	Status    StatusDelivery
	CompanyId uint
	DriverId  uint
	Recipient string
	Address   string
}
