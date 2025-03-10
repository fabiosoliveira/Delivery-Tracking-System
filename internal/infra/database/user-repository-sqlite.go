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

func (ur UserRepositorySqlite) Save(User *domain.User) error {
	_, err := ur.Db.Exec("INSERT INTO Users (name, email, password, user_type) VALUES (?, ?, ?, ?)", User.Name, User.Email, User.Password, User.UserType)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (ur UserRepositorySqlite) FindByEmail(email *string) (*domain.User, error) {
	row := ur.Db.QueryRow("SELECT * FROM Users WHERE email = ?", email)

	var user domain.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}
