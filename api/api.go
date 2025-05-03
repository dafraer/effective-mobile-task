package api

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/dafraer/effective-mobile-task/enrich"
	"github.com/dafraer/effective-mobile-task/store"
	"go.uber.org/zap"
)

type Service struct {
	host     string //e.g. localhost:8080
	logger   *zap.SugaredLogger
	db       store.Storer
	enricher enrich.Enricher
}

// New returns a new service instance
func New(logger *zap.SugaredLogger, db store.Storer, enricher enrich.Enricher) *Service {
	return &Service{
		logger:   logger,
		db:       db,
		enricher: enricher,
	}
}

// Run runs the service
func (s *Service) Run(ctx context.Context, address string) error {
	//Create a new http server
	srv := &http.Server{
		Addr:        address,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	//Four REST routes
	// /get - get users with filters and pagination
	// /delete - delete user by id
	// /update - update user data
	// /add - add user
	http.HandleFunc("/get", s.getHandler)
	http.HandleFunc("/delete", s.deleteHandler)
	http.HandleFunc("/update", s.updateHandler)
	http.HandleFunc("/add", s.addHandler)

	//Create a channel to listen for errors
	ch := make(chan error)

	//Run the server in a separate goroutine
	go func() {
		defer close(ch)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
			return
		}
		ch <- nil
	}()
	s.logger.Infow("Service is running")
	//Wait for the context to be done or for an error to occur and shutdown the server
	select {
	case <-ctx.Done():
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
		err := <-ch
		if err != nil {
			return err
		}
	case err := <-ch:
		return err
	}
	return nil
}

// getResponse is a struct that contains people and cursor to the next page of data
type getResponse struct {
	Cursor int             `json:"cursor"`
	People []*store.Person `json:"people"`
}

// getHandler returns people with specific filters and pagiantion
func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("Request to getHandler")

	//Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Parse query parameters
	params := r.URL.Query()
	s.logger.Debugw("Request to getHandler", "query values", params)

	//Parse limit
	limitStr := params.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		http.Error(w, "limit must be a positive integer", http.StatusBadRequest)
		s.logger.Errorw("Error converting limit to int", "error", err)
		return
	}

	//Parse cursor
	cursorStr := params.Get("cursor")
	cursor := 0
	if cursorStr != "" {
		cursor, err = strconv.Atoi(cursorStr)
		if err != nil || cursor < 0 {
			http.Error(w, "cursor must be an integer larger or equal to 0", http.StatusBadRequest)
			s.logger.Errorw("Error converting cursor to int", "error", err)
			return
		}
	}

	//Parse age
	ageStr := params.Get("age")
	age := 0
	if ageStr != "" {
		age, err = strconv.Atoi(ageStr)
		if err != nil || age < 0 {
			http.Error(w, "age must be a positive integer", http.StatusBadRequest)
			s.logger.Errorw("Error converting age to int", "error", err)
			return
		}
	}

	//Get the people from the database
	people, err := s.db.GetPeople(r.Context(), &store.Params{
		Limit:       limit,
		Cursor:      cursor,
		Name:        params.Get("name"),
		Surname:     params.Get("surname"),
		Patronymic:  params.Get("patronymic"),
		Age:         age,
		Gender:      params.Get("gender"),
		Nationality: params.Get("nationality"),
	})
	if err != nil {
		http.Error(w, "error getting people", http.StatusInternalServerError)
		s.logger.Errorw("Error getting people", "error", err)
		return
	}

	//Write people as a json response
	resp, err := json.Marshal(getResponse{Cursor: cursor + limit, People: people})
	if err != nil {
		http.Error(w, "error marshalling json", http.StatusInternalServerError)
		s.logger.Errorw("Error marshalling json", "error", err)
		return
	}
	s.logger.Debug("Response from getRequest handler", resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// deleteHandler deletes a person by id
func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("Request to deleteHandler")

	//Check if the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Get id from query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	//Convert the id to an int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		s.logger.Errorw("Error converting id to int", "error", err)
		return
	}
	s.logger.Debugw("Request to deleteHandler", "id", id)

	//Delete person from the database
	if err := s.db.DeletePerson(r.Context(), idInt); err != nil {
		http.Error(w, "error deleting person", http.StatusInternalServerError)
		s.logger.Errorw("Error deleting person", "error", err)
	}
}

type updateRequest struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// updateHandler updates non-emprty fields in the updateRequest
func (s *Service) updateHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("Request to updateHandler")

	//Check if the method is PUT or PATCH
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Parse the request body
	var person updateRequest
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}
	s.logger.Debugw("Request to updateHandler", "body", person)

	//Update person
	if err := s.db.UpdatePerson(r.Context(), &store.Person{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
	}); err != nil {
		http.Error(w, "error editing person", http.StatusInternalServerError)
		s.logger.Errorw("Error editing person", "error", err)
	}
}

type addRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

// addHandler enriches person and saves them to the database
func (s *Service) addHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("Request to addHandler")

	//Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Parse the request body
	var person addRequest
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}
	s.logger.Debugw("Request to addHandler", "body", person)

	//Enrich person struct
	p, err := s.enricher.EnrichPerson(r.Context(), person.Name, person.Surname, person.Patronymic)
	if err != nil {
		http.Error(w, "error enriching person", http.StatusInternalServerError)
		s.logger.Errorw("Error enriching person", "error", err)
		return
	}

	//Insert the person into the database
	id, err := s.db.SavePerson(r.Context(), &store.Person{
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         p.Age,
		Gender:      p.Gender,
		Nationality: p.Nationality,
	})
	if err != nil {
		http.Error(w, "error saving person", http.StatusInternalServerError)
		s.logger.Errorw("Error saving person", "error", err)
		return
	}

	//Write the id as a response
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(id)
	s.logger.Debugw("Response from addHandler", "response", resp)
	w.Write(resp)
}
