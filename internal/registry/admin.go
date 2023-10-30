package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/repository"
)

func (r *registry) NewAdminController() controller.AdminEndpointHandler {
	ar := repository.NewAdminRepository(r.db)

	return controller.NewAdminController(ar)
}
