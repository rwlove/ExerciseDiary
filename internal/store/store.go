package store

import "github.com/rwlove/WorkoutDiary/internal/models"

// Store abstracts data access so both the monolith (SQLite) and the
// split frontend (HTTP API client) can satisfy the same interface.
type Store interface {
	SelectEx() ([]models.Exercise, error)
	InsertEx(ex models.Exercise) error
	DeleteEx(id int) error

	SelectSet() ([]models.Set, error)
	// BulkReplaceSetsByDate atomically replaces all sets for a given date.
	BulkReplaceSetsByDate(date string, sets []models.Set) error

	SelectW() ([]models.BodyWeight, error)
	InsertW(w models.BodyWeight) error
	DeleteW(id int) error
}
