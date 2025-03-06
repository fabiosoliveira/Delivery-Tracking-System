package company

type CompanyDto struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type CompanyDataAccessObject interface {
	ListarCompanies() (*[]CompanyDto, error)
}
