package database

import (
	"database/sql"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/user"
)

type UserRepositorySqlite struct {
	Db *sql.DB
}

func NewUserRepositorySqlite(db *sql.DB) UserRepositorySqlite {
	return UserRepositorySqlite{
		Db: db,
	}
}

func (p UserRepositorySqlite) Save(user *user.User) error {
	_, err := p.Db.Exec("INSERT INTO users (name, email, password, type) VALUES (?, ?, ?, ?)", user.Name, user.Email, user.Password, user.Type)
	if err != nil {
		return err
	}

	return nil
}
