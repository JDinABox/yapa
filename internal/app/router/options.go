package router

import (
	"github.com/JDinABox/yapa/internal/app/config"
	"github.com/JDinABox/yapa/internal/app/session"
	"github.com/JDinABox/yapa/internal/sqlc/db"
)

type option func(*Router) error

func WithConfig(c *config.Config) option {
	return func(r *Router) error {
		if c == nil {
			return ErrNil
		}
		r.Config = c
		return nil
	}
}

func WithData(d *db.Queries) option {
	return func(r *Router) error {
		if d == nil {
			return ErrNil
		}
		r.Data = d
		return nil
	}
}

func WithSessionStore(s *session.Store) option {
	return func(r *Router) error {
		if s == nil {
			return ErrNil
		}
		r.SessionStore = s
		return nil
	}
}
