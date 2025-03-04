package database

import (
	"database/sql"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/user"
)

type UserDataAccessObject struct {
	Db *sql.DB
}

func NewUserDataAccessObject(db *sql.DB) UserDataAccessObject {
	return UserDataAccessObject{Db: db}
}

func (u *UserDataAccessObject) ListarContas() (*[]user.UserDto, error) {
	rows, err := u.Db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.UserDto
	for rows.Next() {
		var user user.UserDto
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}
