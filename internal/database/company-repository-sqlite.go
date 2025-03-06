package database

import (
	"database/sql"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/company"
)

type CompanyRepositorySqlite struct {
	Db *sql.DB
}

func NewCompanyRepositorySqlite(db *sql.DB) CompanyRepositorySqlite {
	return CompanyRepositorySqlite{
		Db: db,
	}
}

func (c CompanyRepositorySqlite) Save(company *company.Company) error {
	_, err := c.Db.Exec("INSERT INTO companies (name, email, password) VALUES (?, ?, ?)", company.Name(), company.Email(), company.Password())
	if err != nil {
		return err
	}

	return nil
}

func (c CompanyRepositorySqlite) FindByEmail(email string) (*company.Company, error) {
	row := c.Db.QueryRow("SELECT * FROM companies WHERE email = ?", email)

	var ID int
	var Name, Email, Password string

	err := row.Scan(&ID, &Name, &Email, &Password)
	if err != nil {
		return nil, err
	}

	return company.NewCompany(Name, Email, Password), nil
}
