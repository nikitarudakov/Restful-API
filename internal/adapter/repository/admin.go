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

// AdminRepository provide methods for managing entities for an entire api service
type AdminRepository interface {
	FindUserProfiles(colName string, page int64) ([]model.Profile, error)
}

// NewAdminRepository implicitly links AdminRepository to adminRepository
func NewAdminRepository(db *mongo.Client) AdminRepository {
	return &adminRepository{db: db}
}

// FindUserProfiles find all user profiles and sets pagination based on
// provided page of type int64. Pagination is implemented with
// methods options.Find().SetLimit() and options.Find().SetSkip()
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
