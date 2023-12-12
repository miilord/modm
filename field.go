package modm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DefaultField represents a structure with default fields for MongoDB documents.
type DefaultField struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

// DefaultUpdatedAt sets the default value for the UpdatedAt field.
func (df *DefaultField) DefaultUpdatedAt() {
	df.UpdatedAt = time.Now()
}

// DefaultCreatedAt sets the default value for the CreatedAt field if it's zero.
func (df *DefaultField) DefaultCreatedAt() {
	if df.CreatedAt.IsZero() {
		df.CreatedAt = time.Now()
	}
}

// DefaultID sets the default value for the _id field if it's zero.
func (df *DefaultField) DefaultID() {
	if df.ID.IsZero() {
		df.ID = primitive.NewObjectID()
	}
}

// BeforeInsert is a hook to set default field values before inserting a document.
func (df *DefaultField) BeforeInsert(ctx context.Context) {
	df.DefaultID()
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
