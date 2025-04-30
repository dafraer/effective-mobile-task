package api

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/dafraer/effective-mobile-task/enrich"
	"github.com/dafraer/effective-mobile-task/store"
	"go.uber.org/zap"
)

type Service struct {
	host   string
	logger *zap.SugaredLogger
	db     store.Storer
}

func New(logger *zap.SugaredLogger, db store.Storer) *Service {
	return &Service{
		logger: logger,
		db:     db,
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
	// /edit - edit user
	// /add - add user
	http.HandleFunc("/get", s.getHandler)
	http.HandleFunc("/delete", s.deleteHandler)
	http.HandleFunc("/edit", s.editHandler)
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

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {

}

// deleteHandler deletes a person by id
func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := s.db.DeletePerson(r.Context(), id); err != nil {
		http.Error(w, "error deleting person", http.StatusInternalServerError)
		s.logger.Errorw("Error deleting person", "error", err)
	}
}

func (s *Service) editHandler(w http.ResponseWriter, r *http.Request) {

}

type addRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

// addHandler enriches person and saves it to the database
func (s *Service) addHandler(w http.ResponseWriter, r *http.Request) {
	//Parse the request body
	var person addRequest
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}

	p, err := enrich.EnrichPerson(r.Context(), person.Name, person.Surname, person.Patronymic)
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
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
