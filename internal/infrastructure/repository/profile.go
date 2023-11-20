package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profileDao"
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

// NewProfileRepository implicitly links repository.ProfileRepository to profileRepository
func NewProfileRepository(profileCollection *mongo.Collection) *ProfileRepository {
	return &ProfileRepository{profileCollection: profileCollection}
}

func (pr *ProfileRepository) FindProfileInStorage(profileName *profileDao.ProfileName) (*model.Profile, error) {
	var profile model.Profile

	keyValue := bson.M{"nickname": profileName}

	searchResult := pr.profileCollection.FindOne(context.TODO(), keyValue)

	if err := searchResult.Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (pr *ProfileRepository) InsertProfileToStorage(profileObj *profileDao.ProfileObj) (*profileDao.InsertResult, error) {
	_, err := pr.FindProfileInStorage(&profileDao.ProfileName{Name: profileObj.Nickname})
	if err == nil {
		return nil, errors.New("profile with such nickname already exists")
	}

	now := time.Now().Unix()
	profileObj.CreatedAt = &now
	profileObj.UpdatedAt = &now

	insertOneResult, err := pr.profileCollection.InsertOne(context.TODO(), profileObj)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	insertedID, ok := insertOneResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	insertResult := &model.InsertResult{Id: insertedID, Username: profileObj.Nickname}

	convertedInsertResult := profileDao.MarshalTogRPCInsertResult(insertResult)

	return convertedInsertResult, nil
}

func (pr *ProfileRepository) DeleteProfileFromStorage(profileName *profileDao.ProfileName) error {
	keyValue := bson.M{"nickname": profileName.Name}

	deleteResult, err := pr.profileCollection.DeleteOne(context.TODO(), keyValue)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return errors.New("none profiles were deleted")
	}

	return nil
}

func (pr *ProfileRepository) UpdateProfileInStorage(profileUpdate *profileDao.ProfileUpdate) error {
	keyValue := bson.M{"nickname": &profileDao.ProfileName{Name: profileUpdate.ProfileName.Name}}

	modelUpdate := profileDao.UnmarshalToUpdate(profileUpdate)
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
func (pr *ProfileRepository) ListProfilesFromStorage(page *profileDao.Page) ([]model.Profile, error) {
	var profiles []model.Profile

	keyValue := bson.M{} // no specific key value

	searchOptions := options.Find().SetLimit(objPerPage * page.Num).SetSkip(objPerPage * (page.Num - 1))

	cursor, err := pr.profileCollection.Find(context.TODO(), keyValue, searchOptions)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}
