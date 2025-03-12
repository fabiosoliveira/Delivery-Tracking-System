package domain

func NewDriver(name string, email string, password string) (*User, error) {
	user, err := NewUser(name, email, password, UserTypeDriver)
	if err != nil {
		return nil, err
	}

	return user, nil
}
