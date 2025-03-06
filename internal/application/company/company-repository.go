package company

type CompanyRepository interface {
	Save(company *Company) error
}
