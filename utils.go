package modm

import (
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPointer returns a pointer to the given value.
func GetPointer[T any](value T) *T {
	return &value
}

// DeepCopy creates a deep copy of the source document using JSON encoding and decoding.
func DeepCopy[T comparable](src T) T {
	srcJSON, err := json.Marshal(src)
	if err != nil {
		return *new(T)
	}

	var dest T
	err = json.Unmarshal(srcJSON, &dest)
	if err != nil {
		return *new(T)
	}

	return dest
}

// StructToBSOND converts structure variables to bson.D.
func StructToBSOND(v interface{}) (doc bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// Generate index models from unique and compound index definitions.
// If uniques/indexes is []string{"name"}, means create index "name"
// If uniques/indexes is []string{"name,-age","uid"}, means create compound indexes: name and -age, then create one index: uid
func IndexesToModel(uniques []string, indexes []string) []mongo.IndexModel {
	var indexesModel []mongo.IndexModel

	for _, index := range uniques {
		var keys bson.D
		vv := strings.Split(index, ",")
		for _, field := range vv {
			key, sort := SplitSortField(field)
			keys = append(keys, primitive.E{key, sort})
		}
		indexesModel = append(indexesModel, mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		})
	}

	for _, index := range indexes {
		var keys bson.D
		vv := strings.Split(index, ",")
		for _, field := range vv {
			key, sort := SplitSortField(field)
			keys = append(keys, primitive.E{key, sort})
		}
		indexesModel = append(indexesModel, mongo.IndexModel{
			Keys: keys,
		})
	}

	return indexesModel
}

// SplitSortField handle sort symbol: "+"/"-" in front of field.
// if "+", return sort as 1
// if "-", return sort as -1
func SplitSortField(field string) (key string, sort int32) {
	key = field
	sort = 1

	if field == "" {
		return
	}

	switch field[0] {
	case '+':
		key = strings.TrimPrefix(field, "+")
		sort = 1
	case '-':
		key = strings.TrimPrefix(field, "-")
		sort = -1
	}

	return
}
