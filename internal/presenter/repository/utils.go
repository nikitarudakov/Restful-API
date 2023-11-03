package repository

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func GenerateUpdateObject(p interface{}, tagAlias string) bson.M {
	var fields = make(map[string]interface{})
	filter := bson.M{"$set": fields}

	v := reflect.ValueOf(p)
	vType := v.Type()
	pTypeOf := reflect.TypeOf(p)

	for i := 0; i < v.NumField(); i++ {
		fieldName := vType.Field(i).Name
		fieldVal := v.Field(i)

		fmt.Println(fieldName)
		field, ok := pTypeOf.FieldByName(fieldName)

		if fieldVal.CanInterface() && !fieldVal.IsZero() && ok {
			tag := field.Tag.Get(tagAlias)

			fields[tag] = fieldVal.Interface()
		}
	}

	return filter
}
