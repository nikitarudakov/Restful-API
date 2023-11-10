package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

type UserRepoManager interface {
	Find(u *model.User) (*model.User, error)
	Create(u *model.User) (*InsertResult, error)
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

func (ur *UserRepository) Create(u *model.User) (*InsertResult, error) {
	_, err := ur.Find(u)
	if err == nil {
		return nil, errors.New("user with such username already exists")
	}

	result, err := ur.db.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, fmt.Errorf("error inserting user data: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	return &InsertResult{Id: insertedID, Username: u.Username}, nil
}

func (ur *UserRepository) Delete(authUsername string) error {
	filter := bson.M{"username": authUsername}

	result, err := ur.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return errors.New("user couldn't be deleted as it was not found in db")
	}

	return nil
}

func (ur *UserRepository) UpdateUsername(newUser *model.User, oldVal string) error {
	filter := bson.M{"username": oldVal}
	update := bson.M{"$set": bson.M{
		"username": newUser.Username,
	}}

	result, err := ur.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if result.ModifiedCount < 1 {
		return errors.New("user data was not updated")
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(u *model.User) error {
	filter := bson.M{"username": u.Username}
	update := bson.M{"$set": bson.M{
		"password": u.Password,
	}}

	result, err := ur.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if result.ModifiedCount < 1 {
		return errors.New("user data was not updated")
	}

	return nil
}
