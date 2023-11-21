package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const objPerPage = 5

type ProfileRepository struct {
	profileCollection *mongo.Collection
}

type ProfileRepoController interface {
	FindProfileInStorage(profileName string) (*model.Profile, error)
	InsertProfileToStorage(profile *model.Profile) (*InsertResult, error)
	DeleteProfileFromStorage(profileName string) error
	UpdateProfileInStorage(profile *model.Update, profileName string) error
	ListProfilesFromStorage(page int64) ([]*model.Profile, error)
}

// NewProfileRepository implicitly links repository.ProfileRepository to profileRepository
func NewProfileRepository(profileCollection *mongo.Collection) *ProfileRepository {
	return &ProfileRepository{profileCollection: profileCollection}
}

func (pr *ProfileRepository) FindProfileInStorage(profileName string) (*model.Profile, error) {
	var profile model.Profile

	keyValue := bson.M{"nickname": profileName}

	searchResult := pr.profileCollection.FindOne(context.TODO(), keyValue)

	if err := searchResult.Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (pr *ProfileRepository) InsertProfileToStorage(profile *model.Profile) (*InsertResult, error) {
	_, err := pr.FindProfileInStorage(profile.Nickname)
	if err == nil {
		return nil, errors.New("profile with such nickname already exists")
	}

	now := time.Now().Unix()
	profile.CreatedAt = &now
	profile.UpdatedAt = &now

	insertResult, err := pr.profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	return &InsertResult{Id: insertedID, Username: profile.Nickname}, nil
}

func (pr *ProfileRepository) DeleteProfileFromStorage(profileName string) error {
	keyValue := bson.M{"nickname": profileName}

	deleteResult, err := pr.profileCollection.DeleteOne(context.TODO(), keyValue)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return errors.New("none profiles were deleted")
	}

	return nil
}

func (pr *ProfileRepository) UpdateProfileInStorage(modelUpdate *model.Update, profileName string) error {
	keyValue := bson.M{"nickname": profileName}

	updateProfileObject := generateUpdateObject(*modelUpdate, "bson")

	updateResult, err := pr.profileCollection.UpdateOne(context.TODO(), keyValue, updateProfileObject)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if updateResult.MatchedCount == 0 {
		return errors.New("there is no profile to update")
	}

	log.Trace().
		Str("service", "update profile").
		Str("doc", fmt.Sprintf("%+v", updateProfileObject)).
		Send()

	return nil
}

// ListProfilesFromStorage find all user profiles and sets pagination based on
// provided page of type int64. Pagination is implemented with
// methods options.Find().SetLimit() and options.Find().SetSkip()
func (pr *ProfileRepository) ListProfilesFromStorage(page int64) ([]*model.Profile, error) {
	var profiles []*model.Profile

	keyValue := bson.M{} // no specific key value

	searchOptions := options.Find().SetLimit(objPerPage * page).SetSkip(objPerPage * (page - 1))

	cursor, err := pr.profileCollection.Find(context.TODO(), keyValue, searchOptions)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}
