package company

import "errors"

type CadastrarCompany struct {
	companyRepository CompanyRepository
}

func NewCadastrarCompany(repository CompanyRepository) *CadastrarCompany {
	return &CadastrarCompany{
		companyRepository: repository,
	}
}

func (c *CadastrarCompany) Execute(user *Company) error {
	company, err := c.companyRepository.FindByEmail(user.Email())
	if err != nil {
		return err
	}

	if company != nil {
		return errors.New("company already exists")
	}

	err = c.companyRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}
