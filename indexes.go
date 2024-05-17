package modm

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureIndexes creates unique and non-unique indexes in the collection.
func (r *Repo[T]) EnsureIndexes(ctx context.Context, uniques []string, indexes []string, indexModels ...mongo.IndexModel) error {
	indexesModel := IndexesToModel(uniques, indexes)
	indexesModel = append(indexesModel, indexModels...)
	var err error
	if len(indexesModel) > 0 {
		_, err = r.collection.Indexes().CreateMany(ctx, indexesModel)
	}
	return err
}

// Indexes is an interface for defining unique and non-unique indexes.
type Indexes interface {
	Uniques() []string
	Indexes() []string
	IndexModels() []mongo.IndexModel
}

// EnsureIndexesByModel creates indexes in the collection based on an Indexes interface.
func (r *Repo[T]) EnsureIndexesByModel(ctx context.Context, model Indexes) error {
	indexesModel := IndexesToModel(model.Uniques(), model.Indexes())
	indexesModel = append(indexesModel, model.IndexModels()...)
	var err error
	if len(indexesModel) > 0 {
		_, err = r.collection.Indexes().CreateMany(ctx, indexesModel)
	}
	return err
}
