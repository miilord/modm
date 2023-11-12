package modm

import (
	"context"
	"testing"
)

// TestDefaultFieldHooks verifies that the DefaultField hooks are functioning correctly.
func TestDefaultFieldHooks(t *testing.T) {
	// Create a DefaultField instance
	df := DefaultField{}

	// Test the BeforeInsert hook
	ctx := context.TODO()
	df.BeforeInsert(ctx)
	if df.Id.IsZero() {
		t.Fatalf("BeforeInsert did not set a default ID")
	}
	if df.CreatedAt.IsZero() {
		t.Fatalf("BeforeInsert did not set a default CreatedAt")
	}
	if df.UpdatedAt.IsZero() {
		t.Fatalf("BeforeInsert did not set a default UpdatedAt")
	}

	// Test the BeforeUpdate hook
	df.BeforeUpdate(ctx)
	if df.UpdatedAt.IsZero() {
		t.Fatalf("BeforeUpdate did not set a default UpdatedAt")
	}
}
