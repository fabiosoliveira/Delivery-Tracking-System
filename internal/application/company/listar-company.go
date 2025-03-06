package company

type ListarCompanys struct {
	companyDao CompanyDataAccessObject
}

func NewListarCompanyDao(companyDao CompanyDataAccessObject) *ListarCompanys {
	return &ListarCompanys{
		companyDao: companyDao,
	}
}

func (c *ListarCompanys) Execute(user *Company) error {
	c.companyDao.ListarCompanies()
	return nil
}
