package database

import (
	"database/sql"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/company"
)

type CompanyDataAccessObject struct {
	Db *sql.DB
}

func NewCompanyDataAccessObject(db *sql.DB) CompanyDataAccessObject {
	return CompanyDataAccessObject{Db: db}
}

func (u *CompanyDataAccessObject) ListarCompanies() (*[]company.CompanyDto, error) {
	rows, err := u.Db.Query("SELECT * FROM companies")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []company.CompanyDto
	for rows.Next() {
		var company company.CompanyDto
		if err := rows.Scan(&company.ID, &company.Name, &company.Email, &company.Password); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return &companies, nil
}
