package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db *mongo.Client
}

// NewUserRepository implicitly links repository.UserRepository to userRepository
// which uses mongo.Client as a database
func NewUserRepository(db *mongo.Client) repository.UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Find(u *model.User) (*model.User, error) {
	coll := getCollection(ur.db, u.TableName())

	filter := bson.M{"username": u.Username}

	result := coll.FindOne(context.TODO(), filter)

	if err := result.Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *userRepository) Create(u *model.User) (interface{}, error) {
	coll := getCollection(ur.db, u.TableName())

	_, err := ur.Find(u)
	if err == nil {
		return nil, errors.New("user with such username already exists")
	}

	result, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return result.InsertedID, nil
}

func (ur *userRepository) UpdateUsername(newUser *model.User, oldVal string) error {
	coll := getCollection(ur.db, newUser.TableName())

	filter := bson.M{"username": oldVal}
	update := bson.M{"$set": bson.M{
		"username": newUser.Username,
	}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return nil
}

func (ur *userRepository) UpdatePassword(u *model.User) error {
	coll := getCollection(ur.db, u.TableName())

	filter := bson.M{"username": u.Username}
	update := bson.M{"$set": bson.M{
		"password": u.Password,
	}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return nil
}
