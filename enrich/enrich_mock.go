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
