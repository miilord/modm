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

type TestPapers struct {
	DefaultField `bson:",inline"`
	Text         string `bson:"text,omitempty" json:"text"`
	NumberID     int    `bson:"number_id,omitempty" json:"numberId"`
}

func (u *TestPapers) Uniques() []string {
	return []string{"number_id"}
}

type TestDB struct {
	DoTransaction DoTransactionFunc
	Papers        *Repo[*TestPapers]
	Counters      *Repo[*TestCounter]
}

func TestTransaction(t *testing.T) {
	database, cleanup := setupTestDatabase(t)
	defer cleanup()
	db := TestDB{
		Papers:   NewRepo[*TestPapers](database.Collection("test_papers")),
		Counters: NewRepo[*TestCounter](database.Collection("test_counters")),
	}
	db.DoTransaction = DoTransaction(database.Client())
	ctx := context.TODO()
	err := db.Papers.EnsureIndexesByModel(ctx, &TestPapers{})
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

		paper, err := db.Papers.InsertOne(sessCtx, &TestPapers{Text: "go", NumberID: paperCounter.Count})
		return paper, err
	})
	require.NoError(t, err)
	paper, ok := result.(*TestPapers)
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

		paper, err := db.Papers.InsertOne(sessCtx, &TestPapers{Text: "go", NumberID: paperCounter.Count})
		return paper, err
	})
	require.NoError(t, err)
	paper2, ok := result2.(*TestPapers)
	require.True(t, ok)
	assert.Equal(t, 2, paper2.NumberID)

	_, err = db.Papers.InsertOne(ctx, &TestPapers{Text: "go", NumberID: 2})
	require.True(t, mongo.IsDuplicateKeyError(err))
}
