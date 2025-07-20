package domain

type DriverRepository interface {
	Save(user *Driver) error
	FindByEmail(email *string) (*Driver, error)
	// FindById(id uint) (User, error)
	ListDriversByCompanyId(id int) ([]Driver, error)
}
