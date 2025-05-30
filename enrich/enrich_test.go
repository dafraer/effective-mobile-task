package enrich

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const defaultTimeOut = time.Second * 5

func TestEnrichPerson(t *testing.T) {

	name := "Ivan"
	surname := "Ivanov"
	patronymic := "Ivanovich"
	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}

	enricher := New(&http.Client{Timeout: defaultTimeOut}, logger.Sugar())
	person, err := enricher.EnrichPerson(context.Background(), name, surname, patronymic)

	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, name, person.Name)
	assert.Equal(t, surname, person.Surname)
	assert.Equal(t, patronymic, person.Patronymic)
	assert.NotEmpty(t, person.Age)
	assert.NotEmpty(t, person.Gender)
	assert.NotEmpty(t, person.Nationality)
}

func TestGetAge(t *testing.T) {
	name := "Ivan"

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}
	enricher := New(&http.Client{Timeout: defaultTimeOut}, logger.Sugar())
	age, err := enricher.getAge(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, age)
}

func TestGetGender(t *testing.T) {
	name := "Ivan"

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}
	enricher := New(&http.Client{Timeout: defaultTimeOut}, logger.Sugar())
	gender, err := enricher.getGender(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, gender)
}

func TestGetNationality(t *testing.T) {
	name := "Ivan"
	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}
	enricher := New(&http.Client{Timeout: defaultTimeOut}, logger.Sugar())

	nationality, err := enricher.getNationality(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, nationality)
}
