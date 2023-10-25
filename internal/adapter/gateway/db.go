package gateway

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	db *mongo.Client
}

type Database interface {
	InsertUpdateItem(u *model.User) (interface{}, error)
}

func NewDatabase(db *mongo.Client) Database {
	return &MongoDB{
		db: db,
	}
}

func (m *MongoDB) InsertUpdateItem(u *model.User) (interface{}, error) {
	coll := m.db.Database(config.C.Database.Name).Collection(u.TableName())

	filter := bson.M{"username": u.Username}
	update := bson.M{"$set": u}

	isUpsert := true
	opt := options.UpdateOptions{
		Upsert: &isUpsert,
	}

	updateResult, err := coll.UpdateOne(context.TODO(), filter, update, &opt)
	if err != nil {
		return nil, fmt.Errorf("error inserting/updating item: %w", err)
	}

	return updateResult.UpsertedID, nil
}
