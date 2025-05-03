package store

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storer interface {
	DeletePerson(ctx context.Context, id int) error
	SavePerson(ctx context.Context, person *Person) (int, error)
	UpdatePerson(ctx context.Context, person *Person) error
	GetPeople(ctx context.Context, params *Params) ([]*Person, error)
}

type Store struct {
	db     *sql.DB
	logger *zap.SugaredLogger
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

func New(db *sql.DB, logger *zap.SugaredLogger) Storer {
	return &Store{
		db:     db,
		logger: logger,
	}
}

// Delete deletes person form the database by their ID
func (s *Store) DeletePerson(ctx context.Context, id int) error {
	s.logger.Debugw("DeletePerson called", "id", id)

	_, err := s.db.ExecContext(ctx, "DELETE FROM people WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}

// SavePerson saves a person to the database and returns the ID of the saved person
func (s *Store) SavePerson(ctx context.Context, person *Person) (int, error) {
	s.logger.Debugw("SavePerson called", "person", *person)

	var id int
	err := s.db.QueryRowContext(ctx, "INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID;",
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&id)
	s.logger.Debugw("Saved person", "id", id)
	return id, err
}

// UpdatePerson updates a person in the database
// It only updates non-zero fields
func (s *Store) UpdatePerson(ctx context.Context, person *Person) error {
	s.logger.Debugw("UpdatePerson called", "person", *person)

	_, err := s.db.ExecContext(ctx, `
	UPDATE people 
	SET 
	name = CASE WHEN $1 <> '' THEN $1 ELSE name END,
	surname = CASE WHEN $2 <> '' THEN $2 ELSE surname END,
	patronymic = CASE WHEN $3 <> '' THEN $3 ELSE patronymic END,
	age = CASE WHEN $4 <> 0 THEN $4 ELSE age END,
	gender = CASE WHEN $5 <> '' THEN $5 ELSE gender END,
	nationality = CASE WHEN $6 <> '' THEN $6 ELSE nationality END
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
// It returns slice of people with ID greater than cursor and specified parameters
func (s *Store) GetPeople(ctx context.Context, params *Params) ([]*Person, error) {
	s.logger.Debugw("GetPeople called", "params", *params)

	//Get people from the database
	rows, err := s.db.QueryContext(ctx, `
	SELECT id, name, surname, patronymic, age, gender, nationality FROM people 
	WHERE 
	id > $1 AND
	($3 = '' OR name = $3) AND
	($4 = '' OR surname = $4) AND
	($5 = '' OR patronymic = $5) AND
	($6 = 0 OR age = $6) AND
	($7 = '' OR gender = $7) AND
	($8 = '' OR nationality = $8)
	ORDER BY id LIMIT $2
	`, params.Cursor, params.Limit, params.Name, params.Surname, params.Patronymic, params.Age, params.Gender, params.Nationality)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//Create a slice of people and scan rows
	people := make([]*Person, 0, params.Limit)
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, err
		}
		people = append(people, &p)
	}
	s.logger.Debugw("Received people from the database", "people", people)

	return people, err
}
