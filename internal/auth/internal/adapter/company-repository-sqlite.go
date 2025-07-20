package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/auth/internal/domain"
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
