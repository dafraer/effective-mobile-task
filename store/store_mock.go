package store

import "context"

type MockStore struct {
}

func NewMockStore() Storer {
	return &MockStore{}
}

func (*MockStore) DeletePerson(ctx context.Context, id int) error {
	return nil
}

func (*MockStore) SavePerson(ctx context.Context, person *Person) (int, error) {
	return 1, nil
}

func (*MockStore) UpdatePerson(ctx context.Context, person *Person) error {
	return nil
}

func (*MockStore) GetPeople(ctx context.Context, params *Params) ([]*Person, error) {
	return []*Person{{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         30,
		Gender:      "male",
		Nationality: "russian",
	}}, nil
}
