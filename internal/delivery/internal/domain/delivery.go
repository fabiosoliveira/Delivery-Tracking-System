package domain

type Delivery struct {
	id         uint
	status     statusDelivery
	company_id uint
	driver_id  uint
	// locations  []Location
	recipient string
	address   string
}

func NewDelivery(company_id uint, driver_id uint, recipient string, address string) *Delivery {
	return &Delivery{
		status:     StatusPendente,
		company_id: company_id,
		driver_id:  driver_id,
		// locations:  []Location{},
		recipient: recipient,
		address:   address,
	}
}

func RestoreDelivery(id uint, status uint8, company_id uint, driver_id uint, recipient string, address string) *Delivery {
	return &Delivery{
		id:         id,
		status:     statusDelivery(status),
		company_id: company_id,
		driver_id:  driver_id,
		// locations:  []Location{},
		recipient: recipient,
		address:   address,
	}
}

func (d *Delivery) Id() uint {
	return d.id
}

func (d *Delivery) Status() statusDelivery {
	return d.status
}

func (d *Delivery) Company_id() uint {
	return d.company_id
}

func (d *Delivery) Driver_id() uint {
	return d.driver_id
}

// func (d *Delivery) Locations() []Location {
// 	return d.locations
// }

func (d *Delivery) Recipient() string {
	return d.recipient
}

func (d *Delivery) Address() string {
	return d.address
}

// func (d *Delivery) AddLocation(id uint, latitude float64, longitude float64, timestamp int64) {
// 	location := Location{
// 		Id:        id,
// 		Latitude:  latitude,
// 		Longitude: longitude,
// 		Timestamp: timestamp,
// 	}
// 	d.locations = append(d.locations, location)
// }
