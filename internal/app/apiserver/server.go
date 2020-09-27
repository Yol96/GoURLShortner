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

// shortReq contains request struct fields for createNewShortLink() route
type shortReq struct {
	Address        string `json:"address" validate:"required"`
	ExpirationTime int64  `json:"expiration_time" validate:"min=0"`
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
	return func(w http.ResponseWriter, r *http.Request) {
		var request shortReq
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		//TODO: add validation

		link, err := s.store.User().Create(request.Address, request.ExpirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jStr, _ := json.Marshal(link)
		w.WriteHeader(http.StatusOK)
		w.Write(jStr)
	}
}

// getShortLinkInfo returns new handlerFunc function for "/info" route
func (s *server) getShortLinkInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		link := values.Get("link")

		info, err := s.store.User().Info(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jStr, _ := json.Marshal(info)
		w.WriteHeader(http.StatusOK)
		w.Write(jStr)
	}
}

// redirectToLink returns new handlerFunc function for "/%" route
func (s *server) redirectToLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		url, err := s.store.User().Get(vars["link"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
