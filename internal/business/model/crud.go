package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertResult struct {
	Id       primitive.ObjectID
	Username string
}
