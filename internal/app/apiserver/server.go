package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/Yol96/GoURLShortner/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// server contains server struct with configured router, logger, storage
type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  *store.Store
}

// newServer creates a new configured server
func newServer(store *store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

// http.Handle interface implementation
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter adds routes
func (s *server) configureRouter() {
	//TODO: add middleware
	s.router.HandleFunc("/new", s.createNewShortLink()).Methods("POST")
	s.router.HandleFunc("/info", s.getShortLinkInfo()).Methods("GET")
	s.router.HandleFunc("/{link:[a-zA-Z0-9]{1,11}}", s.redirectToLink()).Methods("GET")
}

// createNewShortLink returns new handlerFunc function for "/new" route
func (s *server) createNewShortLink() http.HandlerFunc {
	type request struct {
		Address        string `json:"address" validate:"required"`
		ExpirationTime int64  `json:"expiration_time" validate:"min=0"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		//TODO: add validation

		link, err := s.store.User().Create(req.Address, req.ExpirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.respond(w, r, http.StatusCreated, link)
	}
}

// getShortLinkInfo returns new handlerFunc function for "/info" route
func (s *server) getShortLinkInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		link := values.Get("link")

		info, err := s.store.User().Info(link)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, info)
	}
}

// redirectToLink returns new handlerFunc function for "/%" route
func (s *server) redirectToLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		url, err := s.store.User().Get(vars["link"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.logger.Infof("Request for redirect: %s", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	s.logger.Infof("Respond for %+v: code:%d, data:%+v", r.URL, statusCode, data)
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
