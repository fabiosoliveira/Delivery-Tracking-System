package database

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/domain"
)

type UserRepositorySqlite struct {
	Db *sql.DB
}

func NewUserRepositorySqlite(db *sql.DB) domain.UserRepository {
	return UserRepositorySqlite{
		Db: db,
	}
}

func (ur UserRepositorySqlite) Save(user domain.User) error {
	var companyId any
	if user.UserType() == domain.UserTypeDriver {
		companyId = user.(*domain.Driver).CompanyId()
	}

	_, err := ur.Db.Exec("INSERT INTO Users (name, email, password, company_id) VALUES (?, ?, ?, ?)", user.Name(), user.Email(), user.Password(), companyId)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (ur UserRepositorySqlite) FindById(userId uint) (domain.User, error) {
	row := ur.Db.QueryRow("SELECT * FROM Users WHERE id = ?", userId)

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

	if companyId != nil {
		user := domain.RestoreDriver(int(id), name, _email, passwordHash, int(companyId.(int64)))
		return user, nil
	}

	user := domain.RestoreCompany(int(id), name, _email, passwordHash)
	return user, nil
}

func (ur UserRepositorySqlite) FindByEmail(email *string) (domain.User, error) {
	row := ur.Db.QueryRow("SELECT * FROM Users WHERE email = ?", email)

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

	if companyId != nil {
		user := domain.RestoreDriver(int(id), name, _email, passwordHash, int(companyId.(int64)))
		return user, nil
	}

	user := domain.RestoreCompany(int(id), name, _email, passwordHash)
	return user, nil
}

func (ur UserRepositorySqlite) ListDriversByCompanyId(id int) ([]domain.User, error) {
	rows, err := ur.Db.Query("SELECT * FROM Users WHERE company_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("error listing drivers: %w", err)
	}
	defer rows.Close()

	var drivers []domain.User
	for rows.Next() {
		var id, companyId int
		var name, _email, passwordHash string
		if err := rows.Scan(&id, &name, &_email, &passwordHash, &companyId); err != nil {
			return nil, fmt.Errorf("error listing drivers: %w", err)
		}

		driver := domain.RestoreDriver(id, name, _email, passwordHash, companyId)

		drivers = append(drivers, driver)
	}

	return drivers, nil
}
