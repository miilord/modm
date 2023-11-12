package modm

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	DefaultField `bson:",inline"`
	Name         string `bson:"name,omitempty" json:"name"`
	Age          uint   `bson:"age,omitempty" json:"age"`
	Bio          string `bson:"-" json:"bio"`
}

func (u *TestUser) AfterFind(ctx context.Context) {
	u.Bio = fmt.Sprintf("%s is %d years old.", u.Name, u.Age)
}

var AfterInsertOK = false

func (u *TestUser) AfterInsert(ctx context.Context) {
	AfterInsertOK = true
}

var AfterUpdateOK = false

func (u *TestUser) AfterUpdate(ctx context.Context) {
	AfterUpdateOK = true
}

func TestDocument(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()
	repo := NewRepo[*TestUser](db.Collection(testColl))

	ctx := context.TODO()
	user, err := repo.InsertOne(ctx, &TestUser{Name: "gooooo", Age: 6})
	require.NoError(t, err)
	require.NotNil(t, user.CreatedAt)
	require.True(t, AfterInsertOK)

	modifiedCount, err := repo.UpdateOne(ctx, &TestUser{Age: 6}, &TestUser{Name: "goooooo", Age: 7})
	require.NoError(t, err)
	require.NotZero(t, modifiedCount)
	require.True(t, AfterUpdateOK)

	u, err := repo.FindOne(ctx, &TestUser{Age: 7})
	require.NoError(t, err)
	assert.NotEqual(t, u.CreatedAt, u.UpdatedAt)
	assert.Equal(t, "goooooo is 7 years old.", u.Bio)
}
