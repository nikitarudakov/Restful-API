package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

type UserRepoManager interface {
	Find(u *model.User) (*model.User, error)
	Create(u *model.User) (interface{}, error)
	Delete(authUsername string) error
	UpdateUsername(newUser *model.User, oldVal string) error
	UpdatePassword(newUser *model.User) error
}

// NewUserRepository implicitly links repository.UserRepository to userRepository
// which uses mongo.Client as a database
func NewUserRepository(db *mongo.Collection) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Find(u *model.User) (*model.User, error) {
	filter := bson.M{"username": u.Username}

	result := ur.db.FindOne(context.TODO(), filter)

	if err := result.Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Create(u *model.User) (interface{}, error) {
	_, err := ur.Find(u)
	if err == nil {
		return nil, errors.New("user with such username already exists")
	}

	result, err := ur.db.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return result.InsertedID, nil
}

func (ur *UserRepository) Delete(authUsername string) error {
	filter := bson.M{"username": authUsername}

	_, err := ur.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateUsername(newUser *model.User, oldVal string) error {
	filter := bson.M{"username": oldVal}
	update := bson.M{"$set": bson.M{
		"username": newUser.Username,
	}}

	_, err := ur.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(u *model.User) error {
	filter := bson.M{"username": u.Username}
	update := bson.M{"$set": bson.M{
		"password": u.Password,
	}}

	_, err := ur.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return nil
}
