package modm

import (
	"context"
)

// EnsureIndexes creates unique and non-unique indexes in the collection.
func (r *Repo[T]) EnsureIndexes(ctx context.Context, uniques []string, indexes []string) error {
	indexesModel := IndexesToModel(uniques, indexes)
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
}

// EnsureIndexesByModel creates indexes in the collection based on an Indexes interface.
func (r *Repo[T]) EnsureIndexesByModel(ctx context.Context, model Indexes) error {
	indexesModel := IndexesToModel(model.Uniques(), model.Indexes())
	var err error
	if len(indexesModel) > 0 {
		_, err = r.collection.Indexes().CreateMany(ctx, indexesModel)
	}
	return err
}
