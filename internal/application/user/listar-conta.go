package user

type ListarConta struct {
	UserDataAccessObject UserDataAccessObject
}

func NewListarConta(repository UserDataAccessObject) *ListarConta {
	return &ListarConta{
		UserDataAccessObject: repository,
	}
}

func (c *ListarConta) Execute(user *User) error {
	c.UserDataAccessObject.ListarContas()
	return nil
}
