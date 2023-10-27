package repository

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db *mongo.Client
}

func NewUserRepository(db *mongo.Client) repository.UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(u *model.User) (interface{}, error) {
	result, err := ur.UpdateMethod(u)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (ur *userRepository) Find(u *model.User) (*model.User, error) {
	if err := ur.FindMethod(u).Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *userRepository) FindMethod(u *model.User) *mongo.SingleResult {
	coll := ur.getCollection(u.TableName())

	filter, _, _ := ur.getQueryParams(u)

	result := coll.FindOne(context.TODO(), filter)

	return result
}

func (ur *userRepository) UpdateMethod(u *model.User) (interface{}, error) {
	coll := ur.getCollection(u.TableName())

	filter, update, opts := ur.getQueryParams(u)

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return result.UpsertedID, nil
}

func (ur *userRepository) getCollection(colName string) *mongo.Collection {
	return ur.db.Database(config.C.Database.Name).Collection(colName)
}

func (ur *userRepository) getQueryParams(u *model.User) (interface{}, interface{}, *options.UpdateOptions) {
	var filter interface{}
	var update interface{}

	filter = bson.M{"username": u.Username}
	update = bson.M{"$set": u}

	opts := options.Update().SetUpsert(true)

	return filter, update, opts
}
