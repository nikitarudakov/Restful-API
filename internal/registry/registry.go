package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/gateway"
)

type registry struct {
	db gateway.Database
}

type Registry interface {
	NewAppConroller() controller.AppController
}

func NewRegistry(db gateway.Database) Registry {
	return &registry{db: db}
}

func (r *registry) NewAppConroller() controller.AppController {
	return controller.AppController{
		User: r.NewUserController(),
	}
}
