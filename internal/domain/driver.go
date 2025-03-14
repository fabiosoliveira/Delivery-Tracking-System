package domain

type Driver struct {
	User
	companyId uint
}

func NewDriver(name string, email string, password string, companyId uint) (User, error) {

	user, err := newUser(name, email, password)
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		User:      user,
		companyId: companyId,
	}

	return driver, nil
}

func RestoreDriver(id int, name string, _email string, passwordHash string, companyId int) User {
	user := restoreUser(id, name, _email, passwordHash)
	return &Driver{
		User:      user,
		companyId: uint(companyId),
	}
}

func (d *Driver) UserType() userType {
	return UserTypeDriver
}

func (d *Driver) CompanyId() uint {
	return d.companyId
}
