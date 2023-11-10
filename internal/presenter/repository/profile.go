package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const objPerPage = 5

type ProfileRepository struct {
	db    *mongo.Collection
	cache *cache.Database
}

type ProfileRepoManager interface {
	Create(p *model.Profile) (*InsertResult, error)
	Find(p *model.Profile) (*model.Profile, error)
	Update(p *model.Update, authUsername string) error
	Delete(authUsername string) error
	ListUserProfiles(page int64) ([]model.Profile, error)
}

// NewProfileRepository implicitly links repository.ProfileRepository to profileRepository
func NewProfileRepository(db *mongo.Collection, cache *cache.Database) *ProfileRepository {
	return &ProfileRepository{db: db, cache: cache}
}

func (pr *ProfileRepository) Find(p *model.Profile) (*model.Profile, error) {
	filter := bson.M{"nickname": p.Nickname}

	result := pr.db.FindOne(context.TODO(), filter)

	if err := result.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (pr *ProfileRepository) Create(p *model.Profile) (*InsertResult, error) {
	_, err := pr.Find(p)
	if err == nil {
		return nil, errors.New("profile with such nickname already exists")
	}

	now := time.Now().Unix()
	p.CreatedAt = &now
	p.UpdatedAt = &now

	result, err := pr.db.InsertOne(context.TODO(), p)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	return &InsertResult{Id: insertedID, Username: p.Nickname}, nil
}

func (pr *ProfileRepository) Delete(authUsername string) error {
	filter := bson.M{"nickname": authUsername}

	deleteResult, err := pr.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return errors.New("none profiles were deleted")
	}

	return nil
}

func (pr *ProfileRepository) Update(modelUpdate *model.Update, authUsername string) error {
	filter := bson.M{"nickname": authUsername}

	update := GenerateUpdateObject(*modelUpdate, "bson")

	result, err := pr.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("there is no profile to update")
	}

	log.Trace().
		Str("service", "update profile").
		Str("doc", fmt.Sprintf("%+v", update)).
		Send()

	return nil
}

// ListUserProfiles find all user profiles and sets pagination based on
// provided page of type int64. Pagination is implemented with
// methods options.Find().SetLimit() and options.Find().SetSkip()
func (pr *ProfileRepository) ListUserProfiles(page int64) ([]model.Profile, error) {
	var results []model.Profile

	opts := options.Find().SetLimit(objPerPage * page).SetSkip(objPerPage * (page - 1))

	log.Info().Msg("HERE 1")

	if err := pr.cache.GetCache("profiles", &results); err == nil {
		return results, nil
	}

	log.Info().Msg("HERE 2")

	cursor, err := pr.db.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if err = pr.cache.SetCache("profiles", &results); err != nil {
		log.Warn().Str("service", "rating caching").Err(err).Send()
	}

	return results, nil
}
