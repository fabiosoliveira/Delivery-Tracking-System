package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/driver/internal/domain"
)

type DriverRepositorySqlite struct {
	Db *sql.DB
}

func NewDriverRepositorySqlite(db *sql.DB) domain.DriverRepository {
	return DriverRepositorySqlite{
		Db: db,
	}
}

func (ur DriverRepositorySqlite) Save(driver *domain.Driver) error {

	_, err := ur.Db.Exec("INSERT INTO Users (name, email, password, company_id) VALUES (?, ?, ?, ?)", driver.Name(), driver.Email(), driver.Password(), driver.CompanyId())
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (ur DriverRepositorySqlite) FindByEmail(email *string) (*domain.Driver, error) {
	row := ur.Db.QueryRow("SELECT * FROM Users WHERE email = ? AND company_id IS NOT NULL", email)

	var id int64
	var companyId int64
	var name, _email, passwordHash string

	err := row.Scan(&id, &name, &_email, &passwordHash, &companyId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding driver: %w", err)
	}

	user := domain.RestoreDriver(int(id), name, _email, passwordHash, int(companyId))
	return user, nil
}

func (ur DriverRepositorySqlite) ListDriversByCompanyId(id int) ([]domain.Driver, error) {
	rows, err := ur.Db.Query("SELECT * FROM Users WHERE company_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("error listing drivers: %w", err)
	}
	defer rows.Close()

	var drivers []domain.Driver
	for rows.Next() {
		var id, companyId int
		var name, _email, passwordHash string
		if err := rows.Scan(&id, &name, &_email, &passwordHash, &companyId); err != nil {
			return nil, fmt.Errorf("error listing drivers: %w", err)
		}

		driver := domain.RestoreDriver(id, name, _email, passwordHash, companyId)

		drivers = append(drivers, *driver)
	}

	return drivers, nil
}
