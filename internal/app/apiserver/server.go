package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Yol96/GoURLShortner/internal/app/model"
	"github.com/Yol96/GoURLShortner/internal/app/store"
	"github.com/gorilla/handlers"
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
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/new", s.createNewShortLink()).Methods("POST")
	api.HandleFunc("/info", s.getShortLinkInfo()).Methods("GET")
	api.HandleFunc("/{link:[a-zA-Z0-9]{1,11}}", s.redirectToLink()).Methods("GET")

	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
}

// createNewShortLink returns new handlerFunc function for "/new" route
func (s *server) createNewShortLink() http.HandlerFunc {
	type request struct {
		Address        string `json:"address" validate:"required,url"`
		ExpirationTime int    `json:"expiration_time" validate:"min=0"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		sl := &model.Link{
			OriginalAddress: req.Address,
			CreatedAt:       time.Now().String(),
			ExpirationTime:  req.ExpirationTime,
		}

		if err := s.store.User().Create(sl); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := createQRCodeImage(sl); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, sl)
	}
}

// getShortLinkInfo returns new handlerFunc function for "/info" route
func (s *server) getShortLinkInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		link := values.Get("link")

		sl := &model.Link{
			ShortLink: link,
		}

		if err := s.store.User().Info(sl); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, sl)
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
