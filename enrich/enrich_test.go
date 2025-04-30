package enrich

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnrichPerson tests the EnrichPerson function
func TestEnrichPerson(t *testing.T) {
	//Test EnrichPerson function
	name := "Ivan"
	surname := "Ivanov"
	patronymic := "Ivanovich"

	person, err := EnrichPerson(context.Background(), name, surname, patronymic)

	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, name, person.Name)
	assert.Equal(t, surname, person.Surname)
	assert.Equal(t, patronymic, person.Patronymic)
	assert.NotEmpty(t, person.Age)
	assert.NotEmpty(t, person.Gender)
	assert.NotEmpty(t, person.Nationality)
}

// TestGetAge tests the getAge function
func TestGetAge(t *testing.T) {
	//Test getAge function
	name := "Ivan"

	age, err := getAge(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, age)
}

// TestGetGender tests the getGender function
func TestGetGender(t *testing.T) {
	name := "Ivan"

	gender, err := getGender(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, gender)
}

// TestGetNationality tests the getNationality function
func TestGetNationality(t *testing.T) {
	name := "Ivan"

	nationality, err := getNationality(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, nationality)
}
