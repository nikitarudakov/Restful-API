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
	coll := ur.db.Database(config.C.Database.Name).Collection(u.TableName())

	filter := bson.M{"hash_id": u.HashID}
	update := bson.M{"$set": u}

	var upsert = true
	opts := options.UpdateOptions{Upsert: &upsert}

	result, err := coll.UpdateOne(context.TODO(), filter, update, &opts)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return result.UpsertedID, nil
}
