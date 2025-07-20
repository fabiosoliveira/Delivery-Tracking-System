package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth/internal/domain"
)

type CompanyRepositorySqlite struct {
	db *sql.DB
}

func NewUserRepositorySqlite(db *sql.DB) domain.CompanyRepository {
	return CompanyRepositorySqlite{
		db: db,
	}
}

func (c CompanyRepositorySqlite) Save(company *domain.Company) error {
	_, err := c.db.Exec("INSERT INTO Users (name, email, password) VALUES (?, ?, ?, ?)", company.Name(), company.Email(), company.Password())
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

// func (c CompanyRepositorySqlite) findById(userId uint) (domain.Company, error) {
// 	row := c.db.QueryRow("SELECT * FROM Users WHERE id = ?", userId)

// 	var id int64
// 	var companyId any
// 	var name, _email, passwordHash string

// 	err := row.Scan(&id, &name, &_email, &passwordHash, &companyId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("error finding user: %w", err)
// 	}

// 	if companyId != nil {
// 		user := RestoreDriver(int(id), name, _email, passwordHash, int(companyId.(int64)))
// 		return user, nil
// 	}

// 	user := RestoreCompany(int(id), name, _email, passwordHash)
// 	return user, nil
// }

func (c CompanyRepositorySqlite) FindByEmail(email *string) (*domain.Company, error) {
	row := c.db.QueryRow("SELECT * FROM Users WHERE email = ? AND company_id IS NULL", email)

	var id int64
	var companyId any
	var name, _email, passwordHash string

	err := row.Scan(&id, &name, &_email, &passwordHash, &companyId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	company := domain.RestoreCompany(int(id), name, _email, passwordHash)
	return company, nil
}

// func (c CompanyRepositorySqlite) listDriversByCompanyId(id int) ([]domain.Company, error) {
// 	rows, err := ur.db.Query("SELECT * FROM Users WHERE company_id = ?", id)
// 	if err != nil {
// 		return nil, fmt.Errorf("error listing drivers: %w", err)
// 	}
// 	defer rows.Close()

// 	var drivers []domain.Company
// 	for rows.Next() {
// 		var id, companyId int
// 		var name, _email, passwordHash string
// 		if err := rows.Scan(&id, &name, &_email, &passwordHash, &companyId); err != nil {
// 			return nil, fmt.Errorf("error listing drivers: %w", err)
// 		}

// 		driver := domain.RestoreDriver(id, name, _email, passwordHash, companyId)

// 		drivers = append(drivers, driver)
// 	}

// 	return drivers, nil
// }
