package repository

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const objPerPage = 5

type AdminRepository struct {
	db *mongo.Database
}

// NewAdminRepository implicitly links AdminRepository to adminRepository
func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindUserProfiles find all user profiles and sets pagination based on
// provided page of type int64. Pagination is implemented with
// methods options.Find().SetLimit() and options.Find().SetSkip()
func (ar *AdminRepository) FindUserProfiles(colName string, page int64) ([]model.Profile, error) {
	coll := ar.db.Collection(colName)

	opts := options.Find().SetLimit(objPerPage * page).SetSkip(objPerPage * (page - 1))

	var results []model.Profile
	cursor, err := coll.Find(context.TODO(), bson.M{}, opts)

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
