// Package migrate provides utilities for moving data between database backends.
package migrate

import (
	"context"
	"fmt"
	"log"

	"github.com/rwlove/WorkoutDiary/internal/store"
)

// Result holds the counts of records copied per table.
type Result struct {
	Exercises int `json:"exercises"`
	Sets      int `json:"sets"`
	Weight    int `json:"weight"`
}

// SQLiteToPostgres copies all data from the SQLite database at sqlitePath into
// the PostgreSQL store pg. Existing data in PostgreSQL is wiped first so the
// operation is idempotent — re-running it produces the same result.
//
// The entire copy is wrapped in a single transaction; if anything fails the
// PostgreSQL database is left unchanged.
func SQLiteToPostgres(sqlitePath string, pg *store.PostgresStore) (Result, error) {
	sqlite := store.NewSQLite(sqlitePath)

	exs, err := sqlite.SelectEx()
	if err != nil {
		return Result{}, fmt.Errorf("read exercises from sqlite: %w", err)
	}
	sets, err := sqlite.SelectSet()
	if err != nil {
		return Result{}, fmt.Errorf("read sets from sqlite: %w", err)
	}
	ws, err := sqlite.SelectW()
	if err != nil {
		return Result{}, fmt.Errorf("read weight from sqlite: %w", err)
	}

	log.Printf("INFO migrate: read %d exercises, %d sets, %d weight records from SQLite",
		len(exs), len(sets), len(ws))

	ctx := context.Background()
	tx, err := pg.Pool().Begin(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Wipe existing PostgreSQL data before re-importing.
	for _, tbl := range []string{"exercises", "sets", "weight"} {
		if _, err := tx.Exec(ctx, "DELETE FROM "+tbl); err != nil {
			return Result{}, fmt.Errorf("clear table %s: %w", tbl, err)
		}
	}

	for _, ex := range exs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO exercises (gr, place, name, descr, image, color, weight, reps, intensity)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			ex.Group, ex.Place, ex.Name, ex.Descr, ex.Image, ex.Color,
			ex.Weight.String(), ex.Reps, ex.Intensity,
		); err != nil {
			return Result{}, fmt.Errorf("insert exercise %q: %w", ex.Name, err)
		}
	}

	for _, s := range sets {
		if _, err := tx.Exec(ctx,
			`INSERT INTO sets (date, name, color, workout_color, weight, reps, intensity)
			 VALUES ($1::date, $2, $3, $4, $5, $6, $7)`,
			s.Date, s.Name, s.Color, s.WorkoutColor,
			s.Weight.String(), s.Reps, s.Intensity,
		); err != nil {
			return Result{}, fmt.Errorf("insert set %q on %s: %w", s.Name, s.Date, err)
		}
	}

	for _, w := range ws {
		if _, err := tx.Exec(ctx,
			"INSERT INTO weight (date, weight) VALUES ($1::date, $2)",
			w.Date, w.Weight.String(),
		); err != nil {
			return Result{}, fmt.Errorf("insert weight on %s: %w", w.Date, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Result{}, fmt.Errorf("commit: %w", err)
	}

	log.Printf("INFO migrate: completed — %d exercises, %d sets, %d weight records",
		len(exs), len(sets), len(ws))

	return Result{
		Exercises: len(exs),
		Sets:      len(sets),
		Weight:    len(ws),
	}, nil
}
