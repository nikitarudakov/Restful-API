package repository

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection(db *mongo.Client, colName string) *mongo.Collection {
	return db.Database(config.C.Database.Name).Collection(colName)
}
