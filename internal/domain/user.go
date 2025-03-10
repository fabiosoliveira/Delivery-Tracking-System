package domain

type UserType uint8

const (
	UserTypeCompany UserType = iota
	UserTypeDriver
)

var userTypes = [2]string{UserTypeCompany: "company", UserTypeDriver: "driver"}

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
	UserType UserType
}

func (u UserType) String() *string {
	return &userTypes[u]
}
