package stats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"fizz-buzz/internal/service"
	_ "modernc.org/sqlite"
)

// Entry represents a counted request and its hit count.
type Entry struct {
	Params service.FizzBuzzParams `json:"params"`
	Hits   int                    `json:"hits"`
}

// Store exposes operations for statistics persistence.
type Store interface {
	Record(ctx context.Context, params service.FizzBuzzParams) error
	Top(ctx context.Context) (Entry, bool, error)
	Close() error
}

// SQLStore keeps FizzBuzz request statistics in a SQL database.
type SQLStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a ready-to-use SQLite-backed statistics store.
func NewSQLiteStore(dsn string) (*SQLStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	store := &SQLStore{db: db}
	if err := store.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

// Record increments the hit counter for a request.
func (s *SQLStore) Record(ctx context.Context, params service.FizzBuzzParams) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO request_stats (int1, int2, limit_value, str1, str2, hits)
		VALUES (?, ?, ?, ?, ?, 1)
		ON CONFLICT(int1, int2, limit_value, str1, str2)
		DO UPDATE SET hits = hits + 1
	`, params.Int1, params.Int2, params.Limit, params.Str1, params.Str2)
	if err != nil {
		return fmt.Errorf("record stats: %w", err)
	}

	return nil
}

// Top returns the most frequently requested parameters.
func (s *SQLStore) Top(ctx context.Context) (Entry, bool, error) {
	var entry Entry
	row := s.db.QueryRowContext(ctx, `
		SELECT int1, int2, limit_value, str1, str2, hits
		FROM request_stats
		ORDER BY hits DESC, id ASC
		LIMIT 1
	`)

	err := row.Scan(
		&entry.Params.Int1,
		&entry.Params.Int2,
		&entry.Params.Limit,
		&entry.Params.Str1,
		&entry.Params.Str2,
		&entry.Hits,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, fmt.Errorf("read top stats: %w", err)
	}

	return entry, true, nil
}

// Close closes the underlying database connection pool.
func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS request_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			int1 INTEGER NOT NULL,
			int2 INTEGER NOT NULL,
			limit_value INTEGER NOT NULL,
			str1 TEXT NOT NULL,
			str2 TEXT NOT NULL,
			hits INTEGER NOT NULL DEFAULT 1,
			UNIQUE(int1, int2, limit_value, str1, str2)
		)
	`)
	if err != nil {
		return fmt.Errorf("migrate stats table: %w", err)
	}

	return nil
}
