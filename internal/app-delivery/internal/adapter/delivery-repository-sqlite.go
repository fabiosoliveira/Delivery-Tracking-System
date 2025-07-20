package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app-delivery/internal/domain"
)

type DeliveryRepositorySqlite struct {
	Db *sql.DB
}

func NewDeliveryRepositorySqlite(db *sql.DB) domain.DeliveryRepository {
	return DeliveryRepositorySqlite{
		Db: db,
	}
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

		delivery := domain.Delivery{uint(id), domain.StatusDelivery(status), uint(company_id), uint(driver_id), recipient, address}

		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}
