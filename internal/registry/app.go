package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

type registry struct {
	db *mongo.Client
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *mongo.Client) Registry {
	return &registry{db: db}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		User: r.NewUserController(),
	}
}
