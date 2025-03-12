package domain

type Company struct {
	User
}

func NewCompany(name string, email string, password string) (*User, error) {
	user, err := NewUser(name, email, password, UserTypeCompany)
	if err != nil {
		return nil, err
	}

	return user, nil
}
