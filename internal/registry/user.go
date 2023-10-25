package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/repository"
)

func (r *registry) NewUserController() controller.User {
	ur := repository.NewUserRepository(r.db)

	return controller.NewUserController(ur)
}
