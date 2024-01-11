package modm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetPointer(t *testing.T) {
	value := 42
	ptr := GetPointer(value)
	assert.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

func TestIsDocumentExists(t *testing.T) {
	// Case 1: Error is nil
	exists, err := IsDocumentExists(nil)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if !exists {
		t.Error("Expected true, got false")
	}

	// Case 2: Error is mongo.ErrNoDocuments
	exists, err = IsDocumentExists(mongo.ErrNoDocuments)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if exists {
		t.Error("Expected false, got true")
	}

	// Case 3: Other error
	otherError := errors.New("some other error")
	exists, err = IsDocumentExists(otherError)
	if err != otherError {
		t.Errorf("Expected %v error, got %v", otherError, err)
	}
	if exists {
		t.Error("Expected false, got true")
	}
}

func TestIndexesToModel(t *testing.T) {
	uniques := []string{"name", "uid"}
	indexes := []string{"name,-age", "email"}

	indexModels := IndexesToModel(uniques, indexes)

	// Verify that the generated index models match the expected results
	expectedIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{primitive.E{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{primitive.E{Key: "uid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				primitive.E{Key: "name", Value: 1},
				primitive.E{Key: "age", Value: -1},
			},
		},
		{
			Keys: bson.D{primitive.E{Key: "email", Value: 1}},
		},
	}

	if len(indexModels) != len(expectedIndexModels) {
		t.Errorf("Expected %d index models, but got %d", len(expectedIndexModels), len(indexModels))
	}

	// Verify each index model
	for i, model := range indexModels {
		if expectedIndexModels[i].Options == nil {
			assert.Nil(t, model.Options)
		} else {
			assert.Equal(t, *model.Options.Unique, *expectedIndexModels[i].Options.Unique)
		}

		bs, err := bson.Marshal(model.Keys)
		assert.Nil(t, err)
		ebs, err := bson.Marshal(expectedIndexModels[i].Keys)
		assert.Nil(t, err)
		assert.Equal(t, bs, ebs)
	}
}

func TestSplitSortField(t *testing.T) {
	t.Run("Empty Field", func(t *testing.T) {
		key, sort := SplitSortField("")
		assert.Equal(t, "", key)
		assert.Equal(t, int32(1), sort)
	})

	t.Run("Sort Ascending", func(t *testing.T) {
		key, sort := SplitSortField("+fieldName")
		assert.Equal(t, "fieldName", key)
		assert.Equal(t, int32(1), sort)
	})

	t.Run("Sort Descending", func(t *testing.T) {
		key, sort := SplitSortField("-fieldName")
		assert.Equal(t, "fieldName", key)
		assert.Equal(t, int32(-1), sort)
	})

	t.Run("Invalid Sort Symbol", func(t *testing.T) {
		key, sort := SplitSortField("*fieldName")
		assert.Equal(t, "*fieldName", key)
		assert.Equal(t, int32(1), sort)
	})
}
