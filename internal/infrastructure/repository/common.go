package repository

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type VoteInsertResult struct {
	Id   primitive.ObjectID
	Vote *model.Vote
}

func generateUpdateObject(update model.Update, tagAlias string) bson.M {
	var fields = make(map[string]interface{})
	filter := bson.M{"$set": fields}

	v := reflect.ValueOf(update)
	vType := v.Type()
	pTypeOf := reflect.TypeOf(update)

	for i := 0; i < v.NumField(); i++ {
		fieldName := vType.Field(i).Name
		fieldVal := v.Field(i)

		field, ok := pTypeOf.FieldByName(fieldName)

		isEmptySlice := fieldVal.Kind() == reflect.Slice && fieldVal.Len() == 0

		if fieldVal.CanInterface() && (!fieldVal.IsZero() || isEmptySlice) && ok {
			tag := field.Tag.Get(tagAlias)

			fields[tag] = fieldVal.Interface()
		}
	}

	fields["updated_at"] = time.Now().Unix()

	return filter
}
