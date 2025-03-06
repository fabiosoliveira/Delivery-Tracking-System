package company

type Company struct {
	name     string
	email    string
	password string
}

func NewCompany(name, email, password string) *Company {
	return &Company{name: name, email: email, password: password}
}

func (c *Company) Name() string {
	return c.name
}

func (c *Company) Email() string {
	return c.email
}

func (c *Company) Password() string {
	return c.password
}
