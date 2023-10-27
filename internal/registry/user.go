package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/repository"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
)

func (r *registry) NewUserController() controller.User {
	uu := usecase.NewUserUsecase(
		repository.NewUserRepository(r.db),
		repository.NewProfileRepository(r.db),
	)

	return controller.NewUserController(uu)
}
