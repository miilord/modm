package modm

import (
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

// IsDocumentExists checks if a MongoDB FindOne/Find operation returned an error indicating
// the absence of documents. It returns true if documents are found, false if
// no documents are found, and any other error encountered during the operation.
// The function is designed to be used in conjunction with MongoDB FindOne/Find queries.
// Example:
//
//	_, err := db.Mongo.Account.FindOne(context.TODO(), filter)
//	exists, err := IsDocumentExists(err)
//	if err != nil {
//	    return err
//	}
//	if !exists {
//	    return fmt.Errorf("Document not found")
//	}
func IsDocumentExists(err error) (bool, error) {
	// if err == mongo.ErrNoDocuments {
	// 	err = nil
	// }
	// return err == nil, err
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
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
