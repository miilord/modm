package modm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne inserts a single document into the collection.
// Hooks: BeforeInsert, AfterInsert
func (r *Repo[T]) InsertOne(ctx context.Context, doc T, opts ...*options.InsertOneOptions) (T, error) {
	doc.BeforeInsert(ctx)
	defer doc.AfterInsert(ctx)
	_, err := r.collection.InsertOne(ctx, doc, opts...)
	if err != nil {
		return *new(T), err
	}
	return doc, nil
}

// InsertMany inserts multiple documents into the collection.
// Hooks: BeforeInsert, AfterInsert
func (r *Repo[T]) InsertMany(ctx context.Context, docs []T, opts ...*options.InsertManyOptions) error {
	var list []interface{}
	for _, doc := range docs {
		doc.BeforeInsert(ctx)
		defer doc.AfterInsert(ctx)
		list = append(list, doc)
	}
	_, err := r.collection.InsertMany(ctx, list, opts...)
	return err
}

// DeleteOne deletes a single document based on the provided filter.
func (r *Repo[T]) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (deletedCount int64, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	res, err := r.collection.DeleteOne(ctx, filter, opts...)
	return res.DeletedCount, err
}

// DeleteMany deletes multiple documents based on the provided filter.
func (r *Repo[T]) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (deletedCount int64, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	res, err := r.collection.DeleteMany(ctx, filter, opts...)
	return res.DeletedCount, err
}

// UpdateByID updates a document by ID with the provided update/document.
// Hooks(document): BeforeUpdate, AfterUpdate
func (r *Repo[T]) UpdateByID(ctx context.Context, id interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error) {
	return r.UpdateOne(ctx, bson.M{"_id": id}, updateOrDoc, opts...)
}

// UpdateOne updates a single document based on the provided filter and update/document.
// Hooks(document): BeforeUpdate, AfterUpdate
func (r *Repo[T]) UpdateOne(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	if doc, ok := updateOrDoc.(T); ok {
		doc.BeforeUpdate(ctx)
		defer doc.AfterUpdate(ctx)
		res, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": doc}, opts...)
		return res.ModifiedCount, err
	}
	res, err := r.collection.UpdateOne(ctx, filter, updateOrDoc, opts...)
	return res.ModifiedCount, err
}

// UpdateMany updates multiple documents based on the provided filter and update/document.
// Hooks(document): BeforeUpdate, AfterUpdate
func (r *Repo[T]) UpdateMany(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.UpdateOptions) (modifiedCount int64, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	if doc, ok := updateOrDoc.(T); ok {
		doc.BeforeUpdate(ctx)
		defer doc.AfterUpdate(ctx)
		res, err := r.collection.UpdateMany(ctx, filter, bson.M{"$set": doc}, opts...)
		return res.ModifiedCount, err
	}
	res, err := r.collection.UpdateMany(ctx, filter, updateOrDoc, opts...)
	return res.ModifiedCount, err
}

// Find retrieves multiple documents based on the provided filter.
// Hooks: AfterFind
func (r *Repo[T]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (docs []T, err error) {
	docs = make([]T, 0)

	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &docs); err != nil {
		return
	}

	for _, doc := range docs {
		doc.AfterFind(ctx)
	}
	return
}

// FindOne retrieves a single document based on the provided filter.
// Hooks: AfterFind
func (r *Repo[T]) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (doc T, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	err = r.collection.FindOne(ctx, filter, opts...).Decode(&doc)
	if err == nil {
		doc.AfterFind(ctx)
	}
	return
}

// FindOneAndDelete retrieves and deletes a single document based on the provided filter.
// Hooks: AfterFind
func (r *Repo[T]) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (doc T, err error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	err = r.collection.FindOneAndDelete(ctx, filter, opts...).Decode(&doc)
	doc.AfterFind(ctx)
	return
}

// FindOneAndUpdate retrieves, updates, and returns a single document based on the provided filter and update/document.
// Hooks: BeforeUpdate(document), AfterUpdate(document), AfterFind
func (r *Repo[T]) FindOneAndUpdate(ctx context.Context, filter interface{}, updateOrDoc interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error) {
	if f, ok := filter.(T); ok {
		filter, _ = StructToBSOND(f)
	}
	opts = append(opts, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if doc, ok := updateOrDoc.(T); ok {
		doc.BeforeUpdate(ctx)
		defer doc.AfterUpdate(ctx)
		err := r.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": doc}, opts...).Decode(&doc)
		return doc, err
	}
	var doc T
	err := r.collection.FindOneAndUpdate(ctx, filter, updateOrDoc, opts...).Decode(&doc)
	doc.AfterFind(ctx)
	return doc, err
}

// [MODM] Get retrieves a single document by ID(ObjectID) from the collection.
func (r *Repo[T]) Get(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) (T, error) {
	return r.FindOne(ctx, bson.M{"_id": id}, opts...)
}
