package router

import (
	"errors"

	"github.com/JDinABox/yapa/internal/app/config"
	"github.com/JDinABox/yapa/internal/app/session"
	"github.com/JDinABox/yapa/internal/ilog"
	"github.com/JDinABox/yapa/internal/sqlc/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrNil = errors.New("nil")
)

type Router struct {
	Config       *config.Config
	Data         *db.Queries
	SessionStore *session.Store
}

func New(opts ...option) *Router {
	r := &Router{}
	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			ilog.Exit(err)
		}
	}
	return r
}

func (router *Router) Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	r.Route("/api", func(r chi.Router) {})

	return r
}
