package store

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Storer interface {
	DeletePerson(ctx context.Context, id int) error
	SavePerson(ctx context.Context, person *Person) (int, error)
	UpdatePerson(ctx context.Context, person *Person) error
	GetPeople(ctx context.Context, params *Params) ([]*Person, error)
}

type Store struct {
	db *sql.DB
}

type Person struct {
	ID          int    `json:"id"`
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

func (s *Store) DeletePerson(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM people WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}

// SavePerson saves a person to the database and returns the ID of the saved person
func (s *Store) SavePerson(ctx context.Context, person *Person) (int, error) {
	var id int
	err := s.db.QueryRowContext(ctx, "INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID;",
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&id)
	return id, err
}

// EditPerson updates a person in the database
// It only updates the fields that are not empty or zero
func (s *Store) UpdatePerson(ctx context.Context, person *Person) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE people 
	SET CASE WHEN $1 <> '' THEN name = $1 ELSE name END,
	SET CASE WHEN $2 <> '' THEN surname = $2 ELSE surname END,
	SET CASE WHEN $3 <> '' THEN patronymic = $1 ELSE patronymic END,
	SET CASE WHEN $4 <> 0 THEN age = $4 ELSE age END,
	SET CASE WHEN $5 <> '' THEN gender = $5 ELSE gender END,
	SET CASE WHEN $6 <> '' THEN nationality = $6 ELSE nationality END
	WHERE id = $7;
	 `, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.ID)
	return err
}

type Params struct {
	Limit       int
	Cursor      int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

// GetPeople retrieves next page of users from the database
// It returns slice of people with ID greater than cursor and specified parameteres
func (s *Store) GetPeople(ctx context.Context, params *Params) ([]*Person, error) {
	//Build query from params
	var people []*Person
	rows, err := s.db.QueryContext(ctx, `
	SELECT id, name, surname, patronymic, age, gender, nationality FROM people 
	WHERE 
	id > $1 AND
	($3 IS "" OR name = $3) AND
	($4 IS "" OR surname = $4) AND
	($5 IS "" OR patronymic = $5) AND
	($6 IS 0 OR age = $6) AND
	($7 IS "" OR gender = $7) AND
	($8 IS "" OR nationality = $8)
	ORDER BY id LIMIT $2
	`, params.Cursor, params.Limit, params.Name, params.Surname, params.Patronymic, params.Age, params.Gender, params.Nationality)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&people)
	return people, err
}
