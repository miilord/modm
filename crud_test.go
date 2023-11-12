package modm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRepo_InsertOne(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	user, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)
	require.NotNil(t, user.CreatedAt)
	require.NotNil(t, user.UpdatedAt)

	var doc TestUser
	err = db.Collection(testColl).FindOne(ctx, bson.M{"age": 6}).Decode(&doc)
	require.NoError(t, err)
	require.NotNil(t, doc.CreatedAt)
	require.NotNil(t, doc.UpdatedAt)
}

func TestRepo_InsertMany(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Create test data
	ctx := context.TODO()
	docs := []*TestUser{
		{Name: "go", Age: 2},
		{Name: "gooooo", Age: 6},
	}

	err := repo.InsertMany(ctx, docs)
	require.NoError(t, err)

	// Validate that documents were inserted
	var count int64
	count, err = db.Collection(testColl).CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(len(docs)), count)

	var doc TestUser
	err = db.Collection(testColl).FindOne(ctx, bson.M{"age": 6}).Decode(&doc)
	require.NoError(t, err)
	require.NotNil(t, doc.CreatedAt)
	require.NotNil(t, doc.UpdatedAt)
}

func TestRepo_DeleteOne(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert a test document
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertOne(ctx, &TestUser{Name: "go", Age: 2})
	require.NoError(t, err)

	// Call the DeleteOne function
	filter := &TestUser{Name: "go"}
	deletedCount, err := repo.DeleteOne(ctx, filter)
	require.NoError(t, err)

	// Validate that one document was deleted
	require.Equal(t, int64(1), deletedCount)

	var count int64
	count, err = db.Collection(testColl).CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestRepo_DeleteMany(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the DeleteMany function
	filter := bson.M{"age": bson.M{"$gte": 2}}
	deletedCount, err := repo.DeleteMany(ctx, filter)
	require.NoError(t, err)

	// Validate that multiple documents were deleted
	require.Equal(t, int64(2), deletedCount)

	var count int64
	count, err = db.Collection(testColl).CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestRepo_UpdateByID(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert a test document
	ctx := context.TODO()
	result, err := repo.InsertOne(ctx, &TestUser{Name: "go", Age: 2})
	require.NoError(t, err)

	// Call the UpdateByID function
	id := result.Id
	update := &TestUser{Age: 3}
	modifiedCount, err := repo.UpdateByID(ctx, id, update)
	require.NoError(t, err)

	// Validate that one document was modified
	require.Equal(t, int64(1), modifiedCount)

	var doc TestUser
	err = db.Collection(testColl).FindOne(ctx, bson.M{"age": 3}).Decode(&doc)
	require.NoError(t, err)
	require.NotEqual(t, doc.UpdatedAt, doc.CreatedAt)
	require.Equal(t, "go", doc.Name)
	require.Equal(t, uint(3), doc.Age)
}

func TestRepo_UpdateOne(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert a test document
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertOne(ctx, &TestUser{Name: "go", Age: 2})
	require.NoError(t, err)

	// Call the UpdateOne function
	filter := &TestUser{Name: "go"}
	update := &TestUser{Age: 3}
	modifiedCount, err := repo.UpdateOne(ctx, filter, update)
	require.NoError(t, err)

	// Validate that one document was modified
	require.Equal(t, int64(1), modifiedCount)

	var doc TestUser
	err = db.Collection(testColl).FindOne(ctx, bson.M{"age": 3}).Decode(&doc)
	require.NoError(t, err)
	require.NotEqual(t, doc.UpdatedAt, doc.CreatedAt)
	require.Equal(t, "go", doc.Name)
	require.Equal(t, uint(3), doc.Age)
}

func TestRepo_UpdateMany(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the UpdateMany function
	filter := bson.M{}
	update := &TestUser{Age: 3}
	modifiedCount, err := repo.UpdateMany(ctx, filter, update)
	require.NoError(t, err)

	// Validate that multiple documents were modified
	require.Equal(t, int64(2), modifiedCount)

	var doc TestUser
	err = db.Collection(testColl).FindOne(ctx, bson.M{"age": 3}).Decode(&doc)
	require.NoError(t, err)
	require.NotEqual(t, doc.UpdatedAt, doc.CreatedAt)
	require.NotZero(t, doc.Name)
}

func TestRepo_Find(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the Find function
	filter := bson.M{"age": bson.M{"$gte": 0}}
	docs, err := repo.Find(ctx, filter)
	require.NoError(t, err)

	// Validate that multiple documents were retrieved
	require.Equal(t, 2, len(docs))

	result, err := repo.Find(ctx, bson.M{"name": "BUG"})
	require.NoError(t, err)
	require.Zero(t, len(result))
	require.NotNil(t, result)
	require.Equal(t, []*TestUser{}, result)
}

func TestRepo_FindOne(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the FindOne function
	doc, err := repo.FindOne(ctx, &TestUser{Name: "go"})
	require.NoError(t, err)

	// Validate that one document was retrieved
	require.NotNil(t, doc)
	require.Equal(t, uint(2), doc.Age)

	result, err := repo.FindOne(ctx, bson.M{"name": "BUG"})
	require.True(t, err == mongo.ErrNoDocuments)
	require.Nil(t, result)
}

func TestRepo_FindOneAndDelete(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the FindOneAndDelete function
	doc, err := repo.FindOneAndDelete(ctx, &TestUser{Name: "go"})
	require.NoError(t, err)
	// Validate that one document was retrieved and deleted
	require.NotNil(t, doc)
	require.Equal(t, uint(2), doc.Age)

	var count int64
	count, err = db.Collection(testColl).CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestRepo_FindOneAndUpdate(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test documents
	ctx := context.TODO()
	_, err := db.Collection(testColl).InsertMany(ctx, []interface{}{
		&TestUser{Name: "go", Age: 2},
		&TestUser{Name: "gooooo", Age: 6},
	})
	require.NoError(t, err)

	// Call the FindOneAndUpdate function
	filter := &TestUser{Name: "go"}
	update := bson.M{"$set": bson.M{"age": 3}}
	doc, err := repo.FindOneAndUpdate(ctx, filter, update)
	require.NoError(t, err)

	// Validate that one document was retrieved and updated
	require.NotNil(t, doc)
	require.Equal(t, uint(3), doc.Age)

	// Call the FindOneAndUpdate function
	doc, err = repo.FindOneAndUpdate(ctx, &TestUser{Name: "go"}, &TestUser{Age: 1})
	require.NoError(t, err)
	require.Equal(t, uint(1), doc.Age)
	require.NotEqual(t, doc.UpdatedAt, doc.CreatedAt)
	require.Equal(t, "go", doc.Name)
}

func TestRepo_Get(t *testing.T) {
	// Create a test Repo instance
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	// Insert test document
	ctx := context.TODO()
	user, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)
	require.NotNil(t, user.CreatedAt)
	require.NotNil(t, user.UpdatedAt)

	// Call the FindOne function
	doc, err := repo.Get(ctx, user.Id)
	require.NoError(t, err)

	// Validate that one document was retrieved
	require.NotNil(t, doc)
	require.Equal(t, uint(6), doc.Age)
}
