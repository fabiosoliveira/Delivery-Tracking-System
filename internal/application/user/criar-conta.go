package user

type CriarConta struct {
	UserRepository UserRepository
}

func NewCriarConta(repository UserRepository) *CriarConta {
	return &CriarConta{
		UserRepository: repository,
	}
}

func (c *CriarConta) Execute(user *User) error {
	c.UserRepository.Save(user)
	return nil
}
