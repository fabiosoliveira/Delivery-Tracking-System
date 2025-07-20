package domain

type DeliveryRepository interface {
	ListDeliveryByCompanyId(id int) ([]Delivery, error)
}
