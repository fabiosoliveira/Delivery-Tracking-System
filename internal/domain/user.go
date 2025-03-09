package domain

type UserType uint8

const (
	UserTypeCompany UserType = iota
	UserTypeDriver
)

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
	UserType UserType
}
