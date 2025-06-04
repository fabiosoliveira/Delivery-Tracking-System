package domain

type DeliveryRepository interface {
	Save(delivery *Delivery) error
	ListDeliveryByCompanyId(id int) ([]Delivery, error)
}
