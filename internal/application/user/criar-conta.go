package user

type CriarConta struct {
	userRepository UserRepository
}

func NewCriarConta(repository UserRepository) *CriarConta {
	return &CriarConta{
		userRepository: repository,
	}
}

func (c *CriarConta) Execute(user *User) error {
	c.userRepository.Save(user)
	return nil
}
