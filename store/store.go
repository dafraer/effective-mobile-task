package store

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storer interface {
	DeletePerson(ctx context.Context, id int) error
	SavePerson(ctx context.Context, person *Person) (int, error)
	UpdatePerson(ctx context.Context, person *Person) error
	GetPeople(ctx context.Context, params *GetParams) ([]*Person, error)
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
func (s *Store) UpdatePerson(ctx context.Context, person *Person) error {
	s.logger.Debugw("UpdatePerson called", "person", *person)

	_, err := s.db.ExecContext(ctx, `
	UPDATE people 
	SET 
	name = $1 ,
	surname = $2,
	patronymic = $3,
	age = $4,
	gender = $5,
	nationality = $6
	WHERE id = $7;
	 `, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.ID)
	return err
}

type GetParams struct {
	Limit       int
	Cursor      *int
	Name        *string
	Surname     *string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationality *string
}

// NewParams populates GetParams struct and returns a pointer to it
func NewParams(limit, cursor, age int, name, surname, patronymic, gender, nationality string) *GetParams {
	params := &GetParams{}
	params.Limit = limit
	params.Cursor = &cursor
	if name != "" {
		params.Name = &name
	}
	if surname != "" {
		params.Surname = &surname
	}
	if patronymic != "" {
		params.Patronymic = &patronymic
	}
	if age != 0 {
		params.Age = &age
	}
	if gender != "" {
		params.Gender = &gender
	}
	if nationality != "" {
		params.Nationality = &nationality
	}
	return params
}

// GetPeople retrieves next page of users from the database
// It returns slice of people with ID greater than cursor and specified parameters
func (s *Store) GetPeople(ctx context.Context, params *GetParams) ([]*Person, error) {
	s.logger.Debugw("GetPeople called", "params", *params)

	//Build query
	q := strings.Builder{}
	paramList := []interface{}{params.Limit, params.Name, params.Surname, params.Patronymic, params.Age, params.Gender, params.Nationality}
	q.WriteString("SELECT id, name, surname, patronymic, age, gender, nationality FROM people WHERE ")
	if params.Cursor != nil {
		q.WriteString("id > $8 AND")
		paramList = append(paramList, params.Cursor)
	}
	q.WriteString(`	($2::TEXT IS NULL OR name = $2) AND
	($3::TEXT IS NULL OR surname = $3) AND
	($4::TEXT IS NULL OR patronymic = $4) AND
	($5::INTEGER IS NULL OR age = $5) AND
	($6::TEXT IS NULL OR gender = $6) AND
	($7::TEXT IS NULL OR nationality = $7)
	ORDER BY id LIMIT $1;
	`)
	//Get people from the database
	rows, err := s.db.QueryContext(ctx, q.String(), paramList...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return people, nil
}
