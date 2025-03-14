package domain

type Company struct {
	User
}

func NewCompany(name string, email string, password string) (User, error) {
	company := &Company{}

	user, err := newUser(name, email, password)
	if err != nil {
		return nil, err
	}

	company.User = user

	return company, nil
}

func RestoreCompany(id int, name string, _email string, passwordHash string) User {
	return &Company{
		User: restoreUser(id, name, _email, passwordHash),
	}
}

func (c *Company) UserType() userType {
	return UserTypeCompany
}
