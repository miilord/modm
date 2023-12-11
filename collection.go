package modm

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection returns the *mongo.Collection
func (r *Repo[T]) Collection() *mongo.Collection {
	return r.collection
}

// Clone creates a copy of the Collection configured with the given CollectionOptions. The specified options are merged with the existing options on the collection, with the specified options taking precedence.
func (r *Repo[T]) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	return r.collection.Clone(opts...)
}

// Name returns the name of the collection.
func (r *Repo[T]) Name() string {
	return r.collection.Name()
}

// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the collection, see the EstimatedDocumentCount method.
// The filter parameter must be a document and can be used to select which documents contribute to the count. It cannot be nil. An empty document (e.g. bson.D{}) should be used to count all documents in the collection. This will result in a full collection scan.
// The opts parameter can be used to specify options for the operation (see the options.CountOptions documentation).
func (r *Repo[T]) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.collection.CountDocuments(ctx, filter, opts...)
}

// EstimatedDocumentCount executes a count command and returns an estimate of the number of documents in the collection using collection metadata.
// The opts parameter can be used to specify options for the operation (see the options.EstimatedDocumentCountOptions documentation).
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/count/.
func (r *Repo[T]) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return r.collection.EstimatedDocumentCount(ctx, opts...)
}

// Distinct executes a distinct command to find the unique values for a specified field in the collection.
// The fieldName parameter specifies the field name for which distinct values should be returned.
// The filter parameter must be a document containing query operators and can be used to select which documents are considered. It cannot be nil. An empty document (e.g. bson.D{}) should be used to select all documents.
// The opts parameter can be used to specify options for the operation (see the options.DistinctOptions documentation).
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/distinct/.
func (r *Repo[T]) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return r.collection.Distinct(ctx, fieldName, filter, opts...)
}

// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
// The pipeline parameter must be an array of documents, each representing an aggregation stage. The pipeline cannot be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the mongo.Pipeline type can be used. See https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages for a list of valid stages in aggregations.
// The opts parameter can be used to specify options for the operation (see the options.AggregateOptions documentation.)
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/aggregate/.
func (r *Repo[T]) Aggregate(ctx context.Context, pipeline interface{}, res interface{}, opts ...*options.AggregateOptions) error {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err == nil {
		err = cursor.All(ctx, res)
	}
	return err
}

// [MODM] Count counts the number of documents in the collection that match the filter.
func (r *Repo[T]) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.CountDocuments(ctx, filter, opts...)
}

// [MODM] EstimatedCount estimates the number of documents in the collection.
func (r *Repo[T]) EstimatedCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return r.EstimatedDocumentCount(ctx, opts...)
}
