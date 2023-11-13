package modm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	testURI  = "mongodb://localhost:27017/test?readPreference=primary&directConnection=true&ssl=false"
	testDB   = "test"
	testColl = "test"
)

func setupTestDatabase(t *testing.T) (*mongo.Database, func()) {
	clientOpts := options.Client().ApplyURI(testURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	require.NoError(t, err)

	err = client.Ping(context.Background(), readpref.Primary())
	require.NoError(t, err)

	db := client.Database(testDB)

	cleanup := func() {
		err = db.Drop(context.Background())
		require.NoError(t, err)

		err = client.Disconnect(context.Background())
		require.NoError(t, err)
	}
	return db, cleanup
}

func TestRepo_Collection(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))
	// Call the Collection function
	collection := repo.Collection()

	// Verify that the collection is not nil
	if collection == nil {
		t.Errorf("Expected a non-nil collection, but got nil")
	}
}

func TestRepo_Clone(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Call the Clone function
	collection, err := repo.Clone()

	// Verify that the returned collection is not nil and no error occurred
	if collection == nil {
		t.Errorf("Expected a non-nil collection, but got nil")
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
}

func TestRepo_Name(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Call the Name function
	name := repo.Name()

	// Verify that the returned name is $testColl
	assert.Equal(t, testColl, name)
}

func TestRepo_CountDocuments(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	_, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)

	// Call the CountDocuments function with appropriate parameters
	count, err := repo.CountDocuments(ctx, &TestUser{Name: "gooooo"})

	// Verify that the returned count is non-negative and no error occurred
	if count < 0 {
		t.Errorf("Expected a non-negative count, but got %d", count)
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.Equal(t, int64(1), count)
}

func TestRepo_EstimatedDocumentCount(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	_, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)

	// Call the EstimatedDocumentCount function with appropriate parameters
	count, err := repo.EstimatedDocumentCount(ctx)

	// Verify that the returned count is non-negative and no error occurred
	if count < 0 {
		t.Errorf("Expected a non-negative count, but got %d", count)
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.Equal(t, int64(1), count)
}

func TestRepo_Distinct(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	users := []*TestUser{
		{
			Name: "go",
			Age:  2,
		},
		{
			Name: "goo",
			Age:  3,
		},
	}
	ctx := context.TODO()
	err := repo.InsertMany(ctx, users)
	require.NoError(t, err)

	// Call the Distinct function with appropriate parameters
	fieldName := "name"
	filter := bson.M{}
	values, err := repo.Distinct(ctx, fieldName, filter)

	// Verify that the returned values are not empty and no error occurred
	if len(values) == 0 {
		t.Errorf("Expected non-empty distinct values, but got an empty slice")
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.ElementsMatch(t, []interface{}{"go", "goo"}, values)

	values2, err := repo.Distinct(ctx, "name", &TestUser{Age: 2})
	if len(values2) == 0 {
		t.Errorf("Expected non-empty distinct values, but got an empty slice")
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.ElementsMatch(t, []interface{}{"go"}, values2)
}

func TestRepo_Aggregate(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	users := []*TestUser{
		{
			Name: "go",
			Age:  2,
		},
		{
			Name: "goo",
			Age:  3,
		},
		{
			Name: "ggo",
			Age:  3,
		},
	}
	ctx := context.TODO()
	err := repo.InsertMany(ctx, users)
	require.NoError(t, err)

	// Call the Aggregate function with appropriate parameters
	pipeline := []bson.D{
		{{"$match", bson.D{{"age", 2}}}},
		{{"$group", bson.D{{"_id", "$name"}, {"count", bson.D{{"$sum", 1}}}}}},
	}
	var result []bson.M
	err = repo.Aggregate(ctx, pipeline, &result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "go", result[0]["_id"])
	assert.Equal(t, int32(1), result[0]["count"])
}

func TestRepo_Count(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	_, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)

	// Call the CountDocuments function with appropriate parameters
	count, err := repo.Count(ctx, &TestUser{Name: "gooooo"})

	// Verify that the returned count is non-negative and no error occurred
	if count < 0 {
		t.Errorf("Expected a non-negative count, but got %d", count)
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.Equal(t, int64(1), count)
}

func TestRepo_EstimatedCount(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	_, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)

	count, err := repo.EstimatedCount(ctx)

	// Verify that the returned count is non-negative and no error occurred
	if count < 0 {
		t.Errorf("Expected a non-negative count, but got %d", count)
	}
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	assert.Equal(t, int64(1), count)
}
