package modm

import (
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

func TestDeepCopy(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	src := User{
		Name: "John",
		Age:  30,
	}
	dest := DeepCopy[User](src)
	assert.Equal(t, src, dest)
}

func TestStructToBSOND(t *testing.T) {
	// Test struct to bson.D conversion
	src := struct {
		Key   string `bson:"key"`
		Value int    `bson:"value"`
	}{
		Key:   "A",
		Value: 123,
	}
	doc, err := StructToBSOND(src)
	assert.Nil(t, err)
	assert.NotNil(t, doc)

	// Verify the field names are taken from BSON tags
	assert.Len(t, doc, 2)
	assert.Equal(t, doc[0].Key, "key")
	assert.Equal(t, doc[1].Key, "value")

	// Test struct with empty fields
	srcEmpty := struct {
	}{}
	docEmpty, errEmpty := StructToBSOND(srcEmpty)
	assert.Nil(t, errEmpty)
	assert.NotNil(t, docEmpty)

	// Test unsupported type
	unsupported := 42
	_, errUnsupported := StructToBSOND(unsupported)
	assert.NotNil(t, errUnsupported)
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
