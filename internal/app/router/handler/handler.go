package handler

import (
	"github.com/JDinABox/yapa/internal/app/config"
	"github.com/JDinABox/yapa/internal/app/session"
	"github.com/JDinABox/yapa/internal/sqlc/db"
)

type Handler struct {
	Config *config.Config
	DB     *db.Queries
	Store  *session.Store
}
