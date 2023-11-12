package modm

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type testIndex struct {
	Key    bson.D
	Name   string
	Unique bool
}

func verifyIndexExists(t *testing.T, iv mongo.IndexView, expected testIndex) {
	cursor, err := iv.List(context.Background())
	assert.Nil(t, err, "List error: %v", err)

	var found bool
	for cursor.Next(context.Background()) {
		var idx testIndex
		err = cursor.Decode(&idx)
		assert.Nil(t, err, "Decode error: %v", err)

		if idx.Name == expected.Name {
			if expected.Key != nil {
				assert.Equal(t, expected.Key, idx.Key, "key document mismatch; expected %v, got %v", expected.Key, idx.Key)
			}
			assert.Equal(t, expected.Unique, idx.Unique)
			found = true
		}
	}
	assert.Nil(t, cursor.Err(), "cursor error: %v", err)
	assert.True(t, found, "expected to find index %v but was not found", expected.Name)
}

func TestRepo_EnsureIndexes(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Define unique and non-unique indexes for testing
	uniques := []string{"name"}
	indexes := []string{"name,-age", "-name", "age,-name"}

	// Call the EnsureIndexes function
	ctx := context.TODO()
	err := repo.EnsureIndexes(ctx, uniques, indexes)
	require.NoError(t, err)

	// Validate that indexes were created
	indexView := repo.collection.Indexes()

	// Validate unique indexes
	for _, item := range IndexesToModel(uniques, indexes) {
		key := item.Keys.(bson.D)
		name := ""
		for i, kv := range key {
			if i >= 1 {
				name += "_"
			}
			name += fmt.Sprintf("%s_%d", kv.Key, kv.Value)
		}
		unique := false
		if item.Options != nil {
			unique = *item.Options.Unique
		}
		verifyIndexExists(t, indexView, testIndex{
			Key:    key,
			Name:   name,
			Unique: unique,
		})
	}
}

func (u *TestUser) Uniques() []string {
	return []string{"name"}
}

func (u *TestUser) Indexes() []string {
	return []string{"name,-age", "-name", "age,-name"}
}

func TestRepo_EnsureIndexesByModel(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Call the EnsureIndexesByModel function
	ctx := context.TODO()
	testUserModel := TestUser{}
	err := repo.EnsureIndexesByModel(ctx, &testUserModel)
	require.NoError(t, err)

	// Validate that indexes were created
	indexView := repo.collection.Indexes()

	// Validate unique indexes
	for _, item := range IndexesToModel(testUserModel.Uniques(), testUserModel.Indexes()) {
		key := item.Keys.(bson.D)
		name := ""
		for i, kv := range key {
			if i >= 1 {
				name += "_"
			}
			name += fmt.Sprintf("%s_%d", kv.Key, kv.Value)
		}
		unique := false
		if item.Options != nil {
			unique = *item.Options.Unique
		}
		verifyIndexExists(t, indexView, testIndex{
			Key:    key,
			Name:   name,
			Unique: unique,
		})
	}
}
