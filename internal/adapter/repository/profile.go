package repository

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ProfileRepository struct {
	db *mongo.Collection
}

// NewProfileRepository implicitly links repository.ProfileRepository to profileRepository
func NewProfileRepository(db *mongo.Collection) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (pr *ProfileRepository) Find(p *model.Profile) (*model.Profile, error) {
	filter := bson.M{"nickname": p.Nickname}

	result := pr.db.FindOne(context.TODO(), filter)

	if err := result.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (pr *ProfileRepository) Create(p *model.Profile) (interface{}, error) {
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

	return result.InsertedID, nil
}

func (pr *ProfileRepository) Update(p *model.Profile) error {
	now := time.Now().Unix()
	p.UpdatedAt = &now

	filter := bson.M{"nickname": p.Nickname}

	update := bson.M{"$set": bson.M{
		"nickname":   p.Nickname,
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"updated_at": time.Now().Unix(),
	}}

	result, err := pr.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating/inserting user data: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("there is no profile to update")
	}

	log.Info().Msg("Profile is updated!")

	return nil
}
