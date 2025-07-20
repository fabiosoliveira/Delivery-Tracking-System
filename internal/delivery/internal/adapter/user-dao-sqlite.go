package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/delivery/internal/domain"
)

type UserDaoSqlite struct {
	Db *sql.DB
}

func NewUserDaoSqlite(db *sql.DB) domain.UserDao {
	return UserDaoSqlite{
		Db: db,
	}
}

// FindCompanyById implements domain.UserDao.
func (ur UserDaoSqlite) FindCompanyById(id uint) (*domain.User, error) {
	row := ur.Db.QueryRow("SELECT name FROM Users WHERE id = ? AND company_id IS NULL", id)

	return ur.toUser(row)
}

// FindDriverById implements domain.UserDao.
func (ur UserDaoSqlite) FindDriverById(id uint) (*domain.User, error) {
	row := ur.Db.QueryRow("SELECT name FROM Users WHERE id = ? AND company_id IS NOT NULL", id)

	return ur.toUser(row)
}

func (ur UserDaoSqlite) toUser(row *sql.Row) (*domain.User, error) {

	var name string

	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	user := domain.User{
		Name: name,
	}
	return &user, nil
}

func (ur UserDaoSqlite) ListDriversByCompanyId(id int) ([]domain.User, error) {
	rows, err := ur.Db.Query("SELECT * FROM Users WHERE company_id = ? AND company_id IS NOT NULL", id)
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

		driver := domain.User{
			ID:    id,
			Name:  name,
			Email: _email,
		}

		drivers = append(drivers, driver)
	}

	return drivers, nil
}
