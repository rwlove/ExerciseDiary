package store

import (
	"github.com/aceberg/ExerciseDiary/internal/db"
	"github.com/aceberg/ExerciseDiary/internal/models"
)

// SQLiteStore implements Store using the local SQLite database.
type SQLiteStore struct {
	path string
}

// NewSQLite constructs a SQLiteStore for the given database file path.
func NewSQLite(path string) *SQLiteStore {
	return &SQLiteStore{path: path}
}

func (s *SQLiteStore) SelectEx() ([]models.Exercise, error) {
	return db.SelectEx(s.path), nil
}

func (s *SQLiteStore) InsertEx(ex models.Exercise) error {
	db.InsertEx(s.path, ex)
	return nil
}

func (s *SQLiteStore) DeleteEx(id int) error {
	db.DeleteEx(s.path, id)
	return nil
}

func (s *SQLiteStore) SelectSet() ([]models.Set, error) {
	return db.SelectSet(s.path), nil
}

func (s *SQLiteStore) BulkReplaceSetsByDate(date string, sets []models.Set) error {
	db.BulkDeleteSetsByDate(s.path, date)
	db.BulkAddSets(s.path, sets)
	return nil
}

func (s *SQLiteStore) SelectW() ([]models.BodyWeight, error) {
	return db.SelectW(s.path), nil
}

func (s *SQLiteStore) InsertW(w models.BodyWeight) error {
	db.InsertW(s.path, w)
	return nil
}

func (s *SQLiteStore) DeleteW(id int) error {
	db.DeleteW(s.path, id)
	return nil
}
