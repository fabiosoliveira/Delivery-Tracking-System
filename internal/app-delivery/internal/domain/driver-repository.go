package domain

type DriverRepository interface {
	Save(user *Driver) error
	FindByEmail(email *string) (*Driver, error)
	FindById(id uint) (*Driver, error)
	ListDriversByCompanyId(id int) ([]Driver, error)
}
