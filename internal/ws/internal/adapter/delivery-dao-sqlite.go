package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/ws/internal/domain"
)

type DeliveryDAOSqlite struct {
	Db *sql.DB
}

func NewDeliveryDAOSqlite(db *sql.DB) domain.DeliveryDAO {
	return DeliveryDAOSqlite{
		Db: db,
	}
}

func (d DeliveryDAOSqlite) UpdateLocation(location *domain.Location, deliveryId uint) error {
	_, err := d.Db.Exec("INSERT INTO Locations (delivery_id, latitude, longitude, timestamp) VALUES (?, ?, ?, ?)", deliveryId, location.Latitude, location.Longitude, location.Timestamp)
	if err != nil {
		return fmt.Errorf("error updating location: %w", err)
	}

	return nil
}
