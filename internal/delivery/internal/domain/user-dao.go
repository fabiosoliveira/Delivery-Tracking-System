package domain

type UserDao interface {
	FindCompanyById(id uint) (*User, error)
	FindDriverById(id uint) (*User, error)
	ListDriversByCompanyId(id int) ([]User, error)
}
