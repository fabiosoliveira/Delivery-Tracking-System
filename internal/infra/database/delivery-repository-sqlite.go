package database

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type DeliveryRepositorySqlite struct {
	Db *sql.DB
}

func NewDeliveryRepositorySqlite(db *sql.DB) domain.DeliveryRepository {
	return DeliveryRepositorySqlite{
		Db: db,
	}
}

func (d DeliveryRepositorySqlite) Save(delivery *domain.Delivery) error {
	_, err := d.Db.Exec("INSERT INTO Deliveries (status, recipient, address, company_id, driver_id) VALUES (?, ?, ?, ?, ?)", delivery.Status(), delivery.Recipient(), delivery.Address(), delivery.Company_id(), delivery.Driver_id())
	if err != nil {
		return fmt.Errorf("error saving delivery: %w", err)
	}

	return nil
}

func (d DeliveryRepositorySqlite) UpdateLocation(location *domain.Location, deliveryId uint) error {
	_, err := d.Db.Exec("INSERT INTO Locations (delivery_id, latitude, longitude, timestamp) VALUES (?, ?, ?, ?)", deliveryId, location.Latitude, location.Longitude, location.Timestamp)
	if err != nil {
		return fmt.Errorf("error updating location: %w", err)
	}

	return nil
}

func (d DeliveryRepositorySqlite) FindById(deliveryId uint) (*domain.Delivery, error) {
	row := d.Db.QueryRow("SELECT * FROM Deliveries WHERE id = ?", deliveryId)

	var id, company_id, driver_id uint
	var status uint8
	var recipient, address string

	err := row.Scan(&id, &status, &recipient, &address, &company_id, &driver_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding Delivery: %w", err)
	}

	delivery := domain.RestoreDelivery(id, status, company_id, driver_id, recipient, address)

	return delivery, nil
}

func (d DeliveryRepositorySqlite) ListDeliveryByCompanyId(id int) ([]domain.Delivery, error) {
	rows, err := d.Db.Query("SELECT * FROM Deliveries WHERE company_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("error listing deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []domain.Delivery
	for rows.Next() {
		var status, id, company_id, driver_id int
		var recipient, address string
		if err := rows.Scan(&id, &status, &recipient, &address, &company_id, &driver_id); err != nil {
			return nil, fmt.Errorf("error listing deliveries: %w", err)
		}

		delivery := domain.RestoreDelivery(uint(id), uint8(status), uint(company_id), uint(driver_id), recipient, address)

		deliveries = append(deliveries, *delivery)
	}

	return deliveries, nil
}

func (d DeliveryRepositorySqlite) ListDeliveryByDriverId(driverId int) ([]domain.Delivery, error) {
	rows, err := d.Db.Query("SELECT * FROM Deliveries WHERE driver_id = ?", driverId)
	if err != nil {
		return nil, fmt.Errorf("error listing deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []domain.Delivery
	for rows.Next() {
		var status, id, company_id, driver_id int
		var recipient, address string
		if err := rows.Scan(&id, &status, &recipient, &address, &company_id, &driver_id); err != nil {
			return nil, fmt.Errorf("error listing deliveries: %w", err)
		}

		delivery := domain.RestoreDelivery(uint(id), uint8(status), uint(company_id), uint(driver_id), recipient, address)

		deliveries = append(deliveries, *delivery)
	}

	return deliveries, nil
}
