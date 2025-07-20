package domain

type UserDao interface {
	// Save(user User) error
	// FindByEmail(email *string) (User, error)
	FindCompanyById(id uint) (*User, error)
	FindDriverById(id uint) (*User, error)
	ListDriversByCompanyId(id int) ([]User, error)
}
