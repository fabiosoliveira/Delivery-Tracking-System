package domain

type UserRepository interface {
	Save(user User) error
	FindByEmail(email *string) (User, error)
	ListDriversByCompanyId(id int) ([]User, error)
}
