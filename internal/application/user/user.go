package user

type UserType uint8

const (
	Empresa UserType = iota
	Motorista
)

type User struct {
	name     string
	email    string
	password string
	userType UserType
}

func NewUser(name, email, password string, userType UserType) *User {
	return &User{name: name, email: email, password: password, userType: userType}
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Type() UserType {
	return u.userType
}
