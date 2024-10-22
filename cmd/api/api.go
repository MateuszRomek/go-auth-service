package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mateuszromek/auth/internal/storage"
	"go.uber.org/zap"
)

type application struct {
	cfg    config
	store  storage.Storage
	logger *zap.SugaredLogger
}

type SessionPayload struct {
	UserId string `json:"user_id"`
}

func (a *application) NewRouter() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/session", func(w http.ResponseWriter, r *http.Request) {

			user, err := a.store.Users.CreateUser(r.Context(), "test@test.com", "test", "test")
			if err != nil {
				a.logger.Error(err)
				WriteJSON(w, http.StatusInternalServerError, err.Error())
			}

			session, err := a.store.Sessions.CreateSession(r.Context(), user.Id)
			if err != nil {
				a.logger.Error(err)
				WriteJSON(w, http.StatusInternalServerError, err.Error())
			}

			WriteJSON(w, http.StatusOK, session)

		})
	})

	return r
}
