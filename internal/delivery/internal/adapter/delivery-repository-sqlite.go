package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/delivery/internal/domain"
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

func (d DeliveryRepositorySqlite) FindLocationsByDeliveryID(deliveryID int) ([]domain.Location, error) {
	query := `SELECT latitude, longitude FROM Locations WHERE delivery_id = ? ORDER BY timestamp ASC`
	rows, err := d.Db.Query(query, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []domain.Location
	for rows.Next() {
		var location domain.Location
		if err := rows.Scan(&location.Latitude, &location.Longitude); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}
