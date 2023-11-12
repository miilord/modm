package modm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DefaultField represents a structure with default fields for MongoDB documents.
type DefaultField struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updatedAt"`
}

// DefaultUpdatedAt sets the default value for the updatedAt field.
func (df *DefaultField) DefaultUpdatedAt() {
	df.UpdatedAt = time.Now()
}

// DefaultCreatedAt sets the default value for the createdAt field if it's zero.
func (df *DefaultField) DefaultCreatedAt() {
	if df.CreatedAt.IsZero() {
		df.CreatedAt = time.Now()
	}
}

// DefaultId sets the default value for the _id field if it's zero.
func (df *DefaultField) DefaultId() {
	if df.Id.IsZero() {
		df.Id = primitive.NewObjectID()
	}
}

// BeforeInsert is a hook to set default field values before inserting a document.
func (df *DefaultField) BeforeInsert(ctx context.Context) {
	df.DefaultId()
	df.DefaultCreatedAt()
	df.DefaultUpdatedAt()
}

// AfterInsert is a hook to handle actions after inserting a document.
func (df *DefaultField) AfterInsert(ctx context.Context) {}

// BeforeUpdate is a hook to set default field values before updating a document.
func (df *DefaultField) BeforeUpdate(ctx context.Context) {
	df.DefaultUpdatedAt()
}

// AfterUpdate is a hook to handle actions after updating a document.
func (df *DefaultField) AfterUpdate(ctx context.Context) {}

// AfterFind is a hook to handle actions after finding a document.
func (df *DefaultField) AfterFind(ctx context.Context) {}

// Uniques returns the unique indexes for the collection.
func (df *DefaultField) Uniques() []string {
	return []string{}
}

// Indexes returns the non-unique indexes for the collection.
func (df *DefaultField) Indexes() []string {
	return []string{}
}
