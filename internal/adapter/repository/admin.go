package repository

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const objPerPage = 5

type adminRepository struct {
	db *mongo.Client
}

type AdminRepository interface {
	FindUserProfiles(colName string, page int64) ([]model.Profile, error)
}

func NewAdminRepository(db *mongo.Client) AdminRepository {
	return &adminRepository{db: db}
}

func (pr *adminRepository) FindUserProfiles(colName string, page int64) ([]model.Profile, error) {
	coll := pr.db.Database(config.C.Database.Name).Collection(colName)

	opts := options.Find().SetLimit(objPerPage * page).SetSkip(objPerPage * (page - 1))

	var results []model.Profile
	cursor, err := coll.Find(context.TODO(), bson.M{}, opts)

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
