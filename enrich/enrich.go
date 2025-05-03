package enrich

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

const (
	ageApiUrl         = "https://api.agify.io/"
	genderApiUrl      = "https://api.genderize.io/"
	nationalityApiUrl = "https://api.nationalize.io/"
)

// Person struct represents person enriched with age, gender and nationality
type Person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type Enricher interface {
	EnrichPerson(ctx context.Context, name, surname, patronymic string) (*Person, error)
}

type defaultEnricher struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) Enricher {
	return &defaultEnricher{logger: logger}
}

// Enrich person enriches person struct with age, gender and nationality from APIs and returns enriched struct
func (e *defaultEnricher) EnrichPerson(ctx context.Context, name, surname, patronymic string) (*Person, error) {
	e.logger.Debugw("EnrichPerson called", "name", name, "surname", surname, "patronymic", patronymic)

	//Get age
	age, err := getAge(ctx, name)
	if err != nil {
		return nil, err
	}
	e.logger.Debugw("Received age from API", "age", age)

	//Get gender
	gender, err := getGender(ctx, name)
	if err != nil {
		return nil, err
	}
	e.logger.Debugw("Received gender from API", "gender", gender)

	//Get nationality
	nationality, err := getNationality(ctx, name)
	if err != nil {
		return nil, err
	}
	e.logger.Debugw("Received nationality from API", "nationality", nationality)

	//Return the enriched person
	return &Person{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}, nil
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// getAge makes a request to the agify API and returns most probable age
func getAge(ctx context.Context, name string) (int, error) {
	//Add parameters to the URL
	u, err := url.Parse(ageApiUrl)
	if err != nil {
		return 0, err
	}
	params := url.Values{}
	params.Add("name", name)
	u.RawQuery = params.Encode()

	//Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, err
	}

	//Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	//Parse the response
	var ageResponse AgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ageResponse); err != nil {
		return 0, err
	}

	//Return the age
	return ageResponse.Age, nil
}

type GenderResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

// getGender makes a request to the genderize API and returns the most probable gender
func getGender(ctx context.Context, name string) (string, error) {
	//Add parameters to the URL
	u, err := url.Parse(genderApiUrl)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("name", name)
	u.RawQuery = params.Encode()

	//Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	//Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//Parse the response
	var genderResponse GenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderResponse); err != nil {
		return "", err
	}

	//Return the gender
	return genderResponse.Gender, nil
}

type NationalityResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}
}

// getNationality makes a request to the nationalize API and returns the most probable nationality
func getNationality(ctx context.Context, name string) (string, error) {
	// Add parameters to the URL
	u, err := url.Parse(nationalityApiUrl)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("name", name)
	u.RawQuery = params.Encode()

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	var nationalityResponse NationalityResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalityResponse); err != nil {
		return "", err
	}

	// Return the nationality
	return nationalityResponse.Country[0].CountryID, nil
}
