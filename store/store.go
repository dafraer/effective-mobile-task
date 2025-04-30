package store

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Storer interface {
	DeletePerson(ctx context.Context, id string) error
	SavePerson(ctx context.Context, person *Person) (string, error)
}

type Store struct {
	db *sql.DB
}

type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func New(db *sql.DB) Storer {
	return &Store{
		db: db,
	}
}

func (s *Store) DeletePerson(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SavePerson(ctx context.Context, person *Person) (string, error) {
	return "", nil
}
