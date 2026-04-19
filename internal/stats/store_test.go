package stats

import (
	"context"
	"testing"
	"time"

	"fizz-buzz/internal/service"
)

func TestSQLStoreTopEmpty(t *testing.T) {
	store := newTestStore(t)

	entry, ok, err := store.Top(context.Background())
	if err != nil {
		t.Fatalf("top returned error: %v", err)
	}
	if ok {
		t.Fatalf("expected no top entry, got %+v", entry)
	}
}

func TestSQLStoreRecordAndTop(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	paramsA := service.FizzBuzzParams{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}
	paramsB := service.FizzBuzzParams{Int1: 2, Int2: 7, Limit: 14, Str1: "foo", Str2: "bar"}

	if err := store.Record(ctx, paramsA); err != nil {
		t.Fatalf("record A #1: %v", err)
	}
	if err := store.Record(ctx, paramsB); err != nil {
		t.Fatalf("record B #1: %v", err)
	}
	if err := store.Record(ctx, paramsA); err != nil {
		t.Fatalf("record A #2: %v", err)
	}

	entry, ok, err := store.Top(ctx)
	if err != nil {
		t.Fatalf("top returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected a top entry")
	}

	if entry.Hits != 2 {
		t.Fatalf("expected hits=2, got %d", entry.Hits)
	}
	if entry.Params != paramsA {
		t.Fatalf("unexpected top params: %+v", entry.Params)
	}
}

func TestSQLStoreRespectsCanceledContext(t *testing.T) {
	store := newTestStore(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	params := service.FizzBuzzParams{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}
	if err := store.Record(ctx, params); err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func newTestStore(t *testing.T) *SQLStore {
	t.Helper()

	store, err := NewSQLiteStoreWithTimeout("file::memory:?cache=shared", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("create sqlite store: %v", err)
	}

	t.Cleanup(func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("close sqlite store: %v", closeErr)
		}
	})

	return store
}
