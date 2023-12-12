package modm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDefaultFieldHooks verifies that the DefaultField hooks are functioning correctly.
func TestDefaultFieldHooks(t *testing.T) {
	// Create a DefaultField instance
	df := DefaultField{}

	// Test the BeforeInsert hook
	ctx := context.TODO()
	df.BeforeInsert(ctx)
	if df.ID.IsZero() {
		t.Fatalf("BeforeInsert did not set a default ID")
	}
	if df.CreatedAt.IsZero() {
		t.Fatalf("BeforeInsert did not set a default CreatedAt")
	}
	if df.UpdatedAt.IsZero() {
		t.Fatalf("BeforeInsert did not set a default UpdatedAt")
	}
	df.AfterInsert(ctx)

	// Test the BeforeUpdate hook
	df.BeforeUpdate(ctx)
	if df.UpdatedAt.IsZero() {
		t.Fatalf("BeforeUpdate did not set a default UpdatedAt")
	}
	df.AfterUpdate(ctx)

	df.AfterFind(ctx)

	uniques := df.Uniques()
	require.Zero(t, len(uniques))
	require.NotNil(t, uniques)
	indexes := df.Indexes()
	require.Zero(t, len(indexes))
	require.NotNil(t, indexes)
}
