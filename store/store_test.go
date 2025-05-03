package store

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSavePerson(t *testing.T) {
	//Initialize the store
	store, err := initStore()
	assert.NoError(t, err)

	//Save the person
	person := Person{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         30,
		Gender:      "male",
		Nationality: "russian",
	}
	id, err := store.SavePerson(context.Background(), &person)
	assert.NoError(t, err)

	//Check that id is not empty
	assert.NotEmpty(t, id)
}

func TestDeletePerson(t *testing.T) {
	//Initialize the store
	store, err := initStore()
	assert.NoError(t, err)

	//Save person to the database
	person := Person{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         30,
		Gender:      "male",
		Nationality: "russian",
	}
	id, err := store.SavePerson(context.Background(), &person)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	//Delete person from the database
	assert.NoError(t, store.DeletePerson(context.Background(), id))
}

func TestGetPeople(t *testing.T) {
	//Initialize store
	store, err := initStore()
	assert.NoError(t, err)

	//Set up people
	people := []*Person{
		{
			ID:          1,
			Name:        "Ivan",
			Surname:     "Petrov",
			Patronymic:  "Sergeevich",
			Age:         35,
			Gender:      "male",
			Nationality: "russian",
		},
		{
			ID:          2,
			Name:        "Maria",
			Surname:     "Kuznetsova",
			Patronymic:  "Andreevna",
			Age:         28,
			Gender:      "female",
			Nationality: "ukrainian",
		},
		{
			ID:          3,
			Name:        "Dmitry",
			Surname:     "Smirnov",
			Patronymic:  "Alexeevich",
			Age:         42,
			Gender:      "male",
			Nationality: "russian",
		},
		{
			ID:          4,
			Name:        "Svetlana",
			Surname:     "Popova",
			Patronymic:  "Ivanovna",
			Age:         22,
			Gender:      "female",
			Nationality: "belarusian",
		},
		{
			ID:          5,
			Name:        "Alexei",
			Surname:     "Vasiliev",
			Patronymic:  "Dmitrievich",
			Age:         50,
			Gender:      "male",
			Nationality: "russian",
		},
		{
			ID:          6,
			Name:        "Elena",
			Surname:     "Ivanova",
			Patronymic:  "",
			Age:         61,
			Gender:      "female",
			Nationality: "kazakh",
		},
		{
			ID:          7,
			Name:        "Sergei",
			Surname:     "Mikhailov",
			Patronymic:  "Nikolaevich",
			Age:         29,
			Gender:      "male",
			Nationality: "russian",
		},
		{
			ID:          8,
			Name:        "Olga",
			Surname:     "Fedorova",
			Patronymic:  "Petrovna",
			Age:         45,
			Gender:      "female",
			Nationality: "ukrainian",
		},
		{
			ID:          9,
			Name:        "Nikolai",
			Surname:     "Morozov",
			Patronymic:  "Ivanovich",
			Age:         61,
			Gender:      "male",
			Nationality: "belarusian",
		},
		{
			ID:          10,
			Name:        "Tatiana",
			Surname:     "Ivanova",
			Patronymic:  "Sergeevna",
			Age:         25,
			Gender:      "female",
			Nationality: "russian",
		},
		{
			ID:          11,
			Name:        "Andrei",
			Surname:     "Novikov",
			Patronymic:  "Vladimirovich",
			Age:         38,
			Gender:      "male",
			Nationality: "georgian",
		},
		{
			ID:          12,
			Name:        "Anna",
			Surname:     "Ivanova",
			Patronymic:  "Alexeevna",
			Age:         29,
			Gender:      "female",
			Nationality: "kazakh",
		},
	}

	//Save people to the database
	for _, p := range people {
		id, err := store.SavePerson(context.Background(), p)
		assert.NoError(t, err)
		assert.Equal(t, id, p.ID)
	}

	//Get everyone from the db
	peopleFromDB, err := store.GetPeople(context.Background(), &Params{Limit: 12})
	assert.NoError(t, err)
	for i, p := range peopleFromDB {
		assert.EqualValues(t, p, people[i])
	}

	//Get first 5 people from the db
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 5})
	assert.NoError(t, err)
	for i, p := range peopleFromDB {
		assert.EqualValues(t, p, people[i])
	}

	//Get next 5 people from the db
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 5, Cursor: 5})
	assert.NoError(t, err)
	for i, p := range peopleFromDB {
		assert.EqualValues(t, p, people[i+5])
	}

	//Filter people with Ivanova surname
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 12, Cursor: 1, Surname: "Ivanova"})
	assert.NoError(t, err)
	assert.EqualValues(t, peopleFromDB[0], people[5])
	assert.EqualValues(t, peopleFromDB[1], people[9])
	assert.EqualValues(t, peopleFromDB[2], people[11])

	//Filter women
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 12, Cursor: 1, Gender: "female"})
	assert.NoError(t, err)
	assert.EqualValues(t, peopleFromDB[0], people[1])
	assert.EqualValues(t, peopleFromDB[1], people[3])
	assert.EqualValues(t, peopleFromDB[2], people[5])
	assert.EqualValues(t, peopleFromDB[3], people[7])
	assert.EqualValues(t, peopleFromDB[4], people[9])
	assert.EqualValues(t, peopleFromDB[5], people[11])

	//Filter people who are 61 years old
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 12, Cursor: 1, Age: 61})
	assert.NoError(t, err)
	assert.EqualValues(t, peopleFromDB[0], people[5])
	assert.EqualValues(t, peopleFromDB[1], people[8])

	//Filter people who are kazakh
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{Limit: 12, Cursor: 1, Nationality: "kazakh"})
	assert.NoError(t, err)
	assert.EqualValues(t, peopleFromDB[0], people[5])
	assert.EqualValues(t, peopleFromDB[1], people[11])

	//Filter Andrei Novikov Vladimirovich 38 y.o. georgian
	peopleFromDB, err = store.GetPeople(context.Background(), &Params{
		Limit:       12,
		Cursor:      1,
		Name:        "Andrei",
		Surname:     "Novikov",
		Patronymic:  "Vladimirovich",
		Age:         38,
		Gender:      "male",
		Nationality: "georgian"})
	assert.NoError(t, err)
	assert.EqualValues(t, peopleFromDB[0], people[10])
}

func TestUpdatePerson(t *testing.T) {
	//Initialize store
	store, err := initStore()
	assert.NoError(t, err)

	person := Person{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         30,
		Gender:      "male",
		Nationality: "russian",
	}

	//Save person to the db
	id, err := store.SavePerson(context.Background(), &person)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	//Update person
	person.Age = 1
	person.Nationality = "american"
	assert.NoError(t, store.UpdatePerson(context.Background(), &person))

	//Check if the person has been updated
	people, err := store.GetPeople(context.Background(), &Params{Limit: 1, Name: person.Name, Surname: person.Surname})
	assert.NoError(t, err)
	assert.EqualValues(t, person, *people[0])
}

// initStore initializes store for tests
func initStore() (Storer, error) {
	//Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	dbConnStr := os.Getenv("DB_URI")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	//Perform migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://../db/migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	//Returns the new storage
	return New(db), nil
}
