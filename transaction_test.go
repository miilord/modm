package modm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestCounter struct {
	DefaultField `bson:",inline"`
	Key          string `bson:"key,omitempty" json:"key"`
	Count        int    `bson:"count,omitempty" json:"count"`
}

type TestPaper struct {
	DefaultField `bson:",inline"`
	Text         string `bson:"text,omitempty" json:"text"`
	NumberID     int    `bson:"number_id,omitempty" json:"numberId"`
}

func (u *TestPaper) Uniques() []string {
	return []string{"number_id"}
}

type TestDB struct {
	DoTransaction DoTransactionFunc
	Papers        *Repo[*TestPaper]
	Counters      *Repo[*TestCounter]
}

func TestDoTransaction(t *testing.T) {
	database, cleanup := setupTestDatabase(t)
	defer cleanup()
	db := TestDB{
		Papers:   NewRepo[*TestPaper](database.Collection("test_papers")),
		Counters: NewRepo[*TestCounter](database.Collection("test_counters")),
	}
	db.DoTransaction = DoTransaction(database.Client())
	ctx := context.TODO()
	err := db.Papers.EnsureIndexesByModel(ctx, &TestPaper{})
	require.NoError(t, err)

	_, err = db.Counters.InsertOne(ctx, &TestCounter{Key: "paper", Count: 0})
	require.NoError(t, err)

	result, err := db.DoTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		paperCounter, err := db.Counters.FindOneAndUpdate(
			sessCtx,
			&TestCounter{Key: "paper"},
			bson.M{"$inc": bson.M{"count": 1}},
		)
		if err != nil {
			return nil, err
		}

		paper, err := db.Papers.InsertOne(sessCtx, &TestPaper{Text: "go", NumberID: paperCounter.Count})
		return paper, err
	})
	require.NoError(t, err)
	paper, ok := result.(*TestPaper)
	require.True(t, ok)
	assert.Equal(t, 1, paper.NumberID)

	result2, err := db.DoTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		paperCounter, err := db.Counters.FindOneAndUpdate(
			sessCtx,
			&TestCounter{Key: "paper"},
			bson.M{"$inc": bson.M{"count": 1}},
		)
		if err != nil {
			return nil, err
		}

		paper, err := db.Papers.InsertOne(sessCtx, &TestPaper{Text: "go", NumberID: paperCounter.Count})
		return paper, err
	})
	require.NoError(t, err)
	paper2, ok := result2.(*TestPaper)
	require.True(t, ok)
	assert.Equal(t, 2, paper2.NumberID)

	_, err = db.Papers.InsertOne(ctx, &TestPaper{Text: "go", NumberID: 2})
	require.True(t, mongo.IsDuplicateKeyError(err))

	t.Run("Error Case - StartSession", func(t *testing.T) {
		mockClient := &mongo.Client{}
		// Mock callback function
		mockCallback := func(sessCtx context.Context) (interface{}, error) {
			return nil, nil
		}

		result, err := DoTransaction(mockClient)(context.Background(), mockCallback)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error Case - WithTransaction", func(t *testing.T) {
		ctx := context.TODO()
		result, err := db.DoTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
			paper, err := db.Papers.FindOne(
				sessCtx,
				&TestPaper{NumberID: 404},
			)
			return paper, err
		})
		require.Error(t, err)
		assert.Nil(t, result)
	})
}
