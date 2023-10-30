package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type profileRepository struct {
	db *mongo.Client
}

// NewProfileRepository implicitly links repository.ProfileRepository to profileRepository
func NewProfileRepository(db *mongo.Client) repository.ProfileRepository {
	return &profileRepository{db: db}
}

func (pr *profileRepository) Find(p *model.Profile) (*model.Profile, error) {
	coll := getCollection(pr.db, p.TableName())

	filter := bson.M{"nickname": p.Nickname}

	result := coll.FindOne(context.TODO(), filter)

	if err := result.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (pr *profileRepository) Create(p *model.Profile) (interface{}, error) {
	coll := getCollection(pr.db, p.TableName())

	_, err := pr.Find(p)
	if err == nil {
		return nil, errors.New("profile with such nickname already exists")
	}

	now := time.Now().Unix()
	p.CreatedAt = &now
	p.UpdatedAt = &now

	result, err := coll.InsertOne(context.TODO(), p)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting user data: %w", err)
	}

	return result.InsertedID, nil
}

func (pr *profileRepository) Update(p *model.Profile) error {
	now := time.Now().Unix()
	p.UpdatedAt = &now

	coll := getCollection(pr.db, p.TableName())

	filter := bson.M{"nickname": p.AuthUsername}

	update := bson.M{"$set": bson.M{
		"nickname":   p.Nickname,
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"updated_at": time.Now().Unix(),
	}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("there is no profile to update")
	}

	log.Info().Msg("Profile is updated!")

	return nil
}
