package modm

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Document represents an interface for common document operations.
type Document interface {
	BeforeInsert(ctx context.Context)
	AfterInsert(ctx context.Context)
	BeforeUpdate(ctx context.Context)
	AfterUpdate(ctx context.Context)
	AfterFind(ctx context.Context)
}

// Repo is a generic repository for working with MongoDB collections.
type Repo[T Document] struct {
	collection *mongo.Collection
}

// NewRepo creates a new repository for the given MongoDB collection.
func NewRepo[T Document](collection *mongo.Collection) *Repo[T] {
	repo := Repo[T]{
		collection: collection,
	}
	return &repo
}

// IRepo represents the interface for MongoDB repository operations.
type IRepo[T Document] interface {
	Aggregate(ctx context.Context, pipeline interface{}, res interface{}, opts ...*options.AggregateOptions) error
	Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error)
	Collection() *mongo.Collection
	Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (deletedCount int64, err error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (deletedCount int64, err error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	EnsureIndexes(ctx context.Context, uniques []string, indexes []string, indexModels ...mongo.IndexModel) error
	EnsureIndexesByModel(ctx context.Context, model Indexes) error
	EstimatedCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (docs []T, err error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (doc T, err error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (doc T, err error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error)
	Get(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) (T, error)
	InsertMany(ctx context.Context, docs []T, opts ...*options.InsertManyOptions) error
	InsertOne(ctx context.Context, doc T, opts ...*options.InsertOneOptions) (T, error)
	Name() string
	UpdateByID(ctx context.Context, id interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error)
	UpdateMany(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error)
	UpdateOne(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error)
}

var _ IRepo[*DefaultField] = NewRepo[*DefaultField](nil)
