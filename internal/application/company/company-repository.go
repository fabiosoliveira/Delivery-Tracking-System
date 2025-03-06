package company

type CompanyRepository interface {
	Save(company *Company) error
	FindByEmail(email string) (*Company, error)
}
