package enrich

import "context"

func NewMockEnricher() Enricher {
	return &MockEnricher{}
}

type MockEnricher struct {
}

func (e *MockEnricher) EnrichPerson(ctx context.Context, name, surname, patronymic string) (*Person, error) {
	return &Person{}, nil
}

func (e *MockEnricher) getNationality(ctx context.Context, name string) (string, error) {
	return "", nil
}

func (e *MockEnricher) getAge(ctx context.Context, name string) (int, error) {
	return 0, nil
}

func (e *MockEnricher) getGender(ctx context.Context, name string) (string, error) {
	return "", nil
}
