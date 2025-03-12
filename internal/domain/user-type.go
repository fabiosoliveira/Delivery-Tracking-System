package domain

type userType uint8

const (
	UserTypeCompany userType = iota
	UserTypeDriver
)

var userTypes = [2]string{UserTypeCompany: "company", UserTypeDriver: "driver"}

func (u userType) String() *string {
	return &userTypes[u]
}
