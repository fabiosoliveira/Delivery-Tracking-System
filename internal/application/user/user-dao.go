package user

type UserDto struct {
	ID       int
	Name     string
	Email    string
	Password string
	UserType uint8
}

type UserDataAccessObject interface {
	ListarContas() (*[]UserDto, error)
}
