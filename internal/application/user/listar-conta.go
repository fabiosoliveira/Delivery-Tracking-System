package user

type ListarConta struct {
	userDao UserDataAccessObject
}

func NewListarConta(userDao UserDataAccessObject) *ListarConta {
	return &ListarConta{
		userDao: userDao,
	}
}

func (c *ListarConta) Execute(user *User) error {
	c.userDao.ListarContas()
	return nil
}
