package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JDinABox/yapa/internal/ilog"
)

// /api/auth/logout
func (a *Auth) Logout_Post(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store.Get(r, a.Config.SessionCookieName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ilog.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
	if err = a.Store.Delete(r, w, session); err != nil {
		ilog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
