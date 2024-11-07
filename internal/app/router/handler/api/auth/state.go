package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JDinABox/yapa/internal/app/router/util"
	"github.com/JDinABox/yapa/internal/ilog"
	"github.com/JDinABox/yapa/internal/sqlc/db"
	"github.com/golang/glog"
	"github.com/google/uuid"
)

type AuthState_Response struct {
	Authenticated bool     `json:"authenticated"`
	User          *db.User `json:"user,omitempty"`
}

// /api/auth/state
func (a *Auth) State_Get(w http.ResponseWriter, r *http.Request) {
	// Default response
	response := AuthState_Response{
		Authenticated: false,
	}

	// Get session from cookie
	session, err := a.Store.Get(r, a.Config.SessionCookieName)
	if err != nil {
		ilog.Error(err)
		util.JSONOutS(w, http.StatusInternalServerError, response)
		return
	}

	// Get authenticated from session
	response.Authenticated, _ = session.Values["authenticated"].(bool)
	if !response.Authenticated {
		util.JSONOutS(w, http.StatusUnauthorized, response)
		return
	}

	// Get user from session
	userID, ok := session.Values["userID"].(string)
	if userID == "" || !ok {
		response.Authenticated = false

		util.JSONOutS(w, http.StatusUnauthorized, response)
		return
	}

	// Log user ID
	if glog.V(1) {
		glog.Info("User ID: ", userID)
	}

	// User id to bytes
	idUUID, err := uuid.Parse(userID)
	if err != nil {
		ilog.Error(err)
		util.JSONOutS(w, http.StatusInternalServerError, util.ErrInternalError)
		return
	}

	// Get user from db
	user, err := a.DB.GetUserByID(r.Context(), idUUID[:])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ilog.Error(err)
		util.JSONOutS(w, http.StatusInternalServerError, util.ErrInternalError)
		return
	}

	// set authenticated if user is found
	if !errors.Is(err, sql.ErrNoRows) {
		response.User = &user
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	util.JSONOut(w, response)
}
