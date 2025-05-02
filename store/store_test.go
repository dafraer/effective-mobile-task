package store

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var examplePerson = Person{
	ID:          1,
	Name:        "Ivan",
	Surname:     "Ivanov",
	Patronymic:  "Ivanovich",
	Age:         30,
	Gender:      "male",
	Nationality: "russian",
}

func TestSavePerson(t *testing.T) {
	store, err := initStore()
	assert.NoError(t, err)
	id, err := store.SavePerson(context.Background(), &examplePerson)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestDeletePerson(t *testing.T) {
	store, err := initStore()
	assert.NoError(t, err)
	id, err := store.SavePerson(context.Background(), &examplePerson)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.NoError(t, store.DeletePerson(context.Background(), id))
}

func TestGetPeople(t *testing.T) {

}

func TestUpdatePerson(t *testing.T) {

}

func initStore() (Storer, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	dbConnStr := os.Getenv("DB_URI")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	m.Up()
	return New(db), nil
}
