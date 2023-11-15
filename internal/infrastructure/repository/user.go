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
	userRepo *mongo.Collection
}

type UserRepoController interface {
	FindUserInStorage(username string) (*model.User, error)
	InsertUserToStorage(user *model.User) (*InsertResult, error)
	DeleteUserFromStorage(username string) error
	UpdateUsernameInStorage(userWithNewUsername *model.User, toReplaceUsername string) error
	UpdatePasswordInStorage(u *model.User) error
}

// NewUserRepository implicitly links repository.UserRepository to userRepository
// which uses mongo.Client as a database
func NewUserRepository(userRepo *mongo.Collection) *UserRepository {
	return &UserRepository{userRepo: userRepo}
}

func (ur *UserRepository) FindUserInStorage(username string) (*model.User, error) {
	var user model.User

	keyValue := bson.M{"username": username}

	searchResult := ur.userRepo.FindOne(context.TODO(), keyValue)

	if err := searchResult.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) InsertUserToStorage(user *model.User) (*InsertResult, error) {
	_, err := ur.FindUserInStorage(user.Username)
	if err == nil {
		return nil, errors.New("user with such username already exists")
	}

	insertResult, err := ur.userRepo.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("error inserting user data: %w", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	return &InsertResult{Id: insertedID, Username: user.Username}, nil
}

func (ur *UserRepository) DeleteUserFromStorage(username string) error {
	keyValue := bson.M{"username": username}

	deleteResult, err := ur.userRepo.DeleteOne(context.TODO(), keyValue)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount < 1 {
		return errors.New("user couldn't be deleted as it was not found in db")
	}

	return nil
}

func (ur *UserRepository) UpdateUsernameInStorage(userWithNewUsername *model.User,
	toReplaceUsername string) error {
	keyValue := bson.M{"username": toReplaceUsername}
	updateUserObject := bson.M{"$set": bson.M{
		"username": userWithNewUsername.Username,
	}}

	updateResult, err := ur.userRepo.UpdateOne(context.TODO(), keyValue, updateUserObject)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if updateResult.ModifiedCount < 1 {
		return errors.New("user data was not updated")
	}

	return nil
}

func (ur *UserRepository) UpdatePasswordInStorage(user *model.User) error {
	keyValue := bson.M{"username": user.Username}
	updateUserObject := bson.M{"$set": bson.M{
		"password": user.Password,
	}}

	updateResult, err := ur.userRepo.UpdateOne(context.TODO(), keyValue, updateUserObject)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if updateResult.ModifiedCount < 1 {
		return errors.New("user data was not updated")
	}

	return nil
}
