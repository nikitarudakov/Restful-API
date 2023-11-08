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

type VoteRepository struct {
	db *mongo.Collection
}

type VoteRepoManager interface {
	Find(v *model.Vote, isTarget bool, isSender bool) (*model.Vote, error)
	Create(v *model.Vote) (*VoteInsertResult, error)
	Delete(v *model.Vote, isTarget bool, isSender bool) error
	CalcTotalRating(target string) (*model.Rating, error)
}

func NewVoteRepository(db *mongo.Collection) *VoteRepository {
	return &VoteRepository{db: db}
}

func getFilterForVote(v *model.Vote, isTarget bool, isSender bool) map[string]interface{} {
	var filter = make(map[string]interface{})

	if isTarget {
		filter["target"] = v.Target
	}
	if isSender {
		filter["sender"] = v.Sender
	}

	return filter
}

func (vr *VoteRepository) Find(v *model.Vote, isTarget bool, isSender bool) (*model.Vote, error) {
	var vFromDB = &model.Vote{}

	filter := getFilterForVote(v, isTarget, isSender)

	result := vr.db.FindOne(context.TODO(), filter)

	if err := result.Decode(vFromDB); err != nil {
		return nil, err
	}

	return vFromDB, nil
}

func (vr *VoteRepository) Create(v *model.Vote) (*VoteInsertResult, error) {
	result, err := vr.db.InsertOne(context.TODO(), v)
	if err != nil {
		return nil, fmt.Errorf("error updating/inserting vote data: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("conversion error")
	}

	return &VoteInsertResult{Id: insertedID, Vote: v}, nil
}

func (vr *VoteRepository) Delete(v *model.Vote, isTarget bool, isSender bool) error {
	filter := getFilterForVote(v, isTarget, isSender)

	deleteResult, err := vr.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return errors.New("none votes were deleted")
	}

	return nil
}

func (vr *VoteRepository) CalcTotalRating(target string) (*model.Rating, error) {
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"target", target}}}},
		bson.D{{"$group", bson.D{
			{"_id", "$target"},
			{"rating", bson.D{{"$sum", "$vote"}}},
		}}}}

	cursor, err := vr.db.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []model.Rating
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return &results[0], nil
	}

	return nil, nil
}
