package company

type CadastrarCompany struct {
	companyRepository CompanyRepository
}

func NewCadastrarCompany(repository CompanyRepository) *CadastrarCompany {
	return &CadastrarCompany{
		companyRepository: repository,
	}
}

func (c *CadastrarCompany) Execute(user *Company) error {
	err := c.companyRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}
