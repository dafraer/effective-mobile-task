package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dafraer/effective-mobile-task/enrich"
	"github.com/dafraer/effective-mobile-task/store"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	sugar := logger.Sugar()

	//Create new server for testing
	service := New(sugar, store.NewMockStore(), enrich.NewMockEnricher())

	//Create test server
	server := httptest.NewServer(http.HandlerFunc(service.getHandler))
	srvUrl := server.URL + "?"

	//Make a POST request to make sure it does not work
	params := url.Values{}
	params.Add("limit", "1")
	server.URL = srvUrl + params.Encode()
	req, err := http.NewRequest(http.MethodPost, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 405
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, fmt.Sprintf("expected 405 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a GET request with no limit parameter to make sure it doesn't work
	server.URL = srvUrl
	req, err = http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 400
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected 400 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a GET request with negative  limit parameter to make sure it doesn't work
	params = url.Values{}
	params.Add("limit", "-12")
	server.URL = srvUrl + params.Encode()
	req, err = http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 401
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected 401 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a correct GET request
	//Set the query values
	params = url.Values{}
	params.Add("limit", "5")
	params.Add("cursor", "0")
	params.Add("age", "15")
	params.Add("name", "Ivan")
	params.Add("surname", "Ivanov")
	params.Add("patronymic", "Ivanovich")
	params.Add("gender", "male")
	params.Add("nationality", "russian")

	//Add query values to the URL string
	server.URL = srvUrl + params.Encode()

	//Make a request
	req, err = http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 200
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))
	defer resp.Body.Close()

	//Decode json response into response struct
	var response getResponse
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	//Check that the response is correct
	person := store.Person{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         30,
		Gender:      "male",
		Nationality: "russian",
	}
	assert.Equal(t, *response.NextCursor, 1)
	assert.EqualValues(t, person, *response.People[0])

	//Close response body
	assert.NoError(t, resp.Body.Close())
}

func TestAddHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	sugar := logger.Sugar()

	//Create new server for testing
	service := New(sugar, store.NewMockStore(), enrich.NewMockEnricher())

	//Create test server
	server := httptest.NewServer(http.HandlerFunc(service.addHandler))

	//Make a GET request to make sure it does not work
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 405
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, fmt.Sprintf("expected 405 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a correct POST request
	requestBody := addRequest{
		Name:       "Ivan",
		Surname:    "Ivanov",
		Patronymic: "Ivanovich",
	}
	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)
	req, err = http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 200
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Check that response has correct id
	var id int
	err = json.NewDecoder(resp.Body).Decode(&id)
	assert.Equal(t, 1, id)
	//Close response body
	assert.NoError(t, resp.Body.Close())
}

func TestUpdateHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	sugar := logger.Sugar()

	//Create new server for testing
	service := New(sugar, store.NewMockStore(), enrich.NewMockEnricher())

	//Create test server
	server := httptest.NewServer(http.HandlerFunc(service.updateHandler))

	//Make a GET request to make sure it does not work
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 405
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, fmt.Sprintf("expected 405 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a correct PUT request
	requestBody := updateRequest{
		ID:          1,
		Name:        "Ivan",
		Surname:     "Ivanov",
		Patronymic:  "Ivanovich",
		Age:         14,
		Gender:      "male",
		Nationality: "russian",
	}
	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)
	req, err = http.NewRequest(http.MethodPut, server.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 200
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())
}

func TestDeleteHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	sugar := logger.Sugar()

	//Create new server for testing
	service := New(sugar, store.NewMockStore(), enrich.NewMockEnricher())

	//Create test server
	server := httptest.NewServer(http.HandlerFunc(service.deleteHandler))
	server.URL += "?id=1"

	//Make a GET request to make sure it does not work
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 405
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, fmt.Sprintf("expected 405 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())

	//Make a correct DELETE request
	req, err = http.NewRequest(http.MethodDelete, server.URL, http.NoBody)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is 200
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())
}
