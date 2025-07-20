package domain

type CompanyRepository interface {
	Save(company *Company) error
	FindByEmail(email *string) (*Company, error)
	// findById(id uint) (Company, error)
	// listDriversByCompanyId(id int) ([]Company, error)
}
