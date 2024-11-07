package util

import (
	"fmt"
	"net/http"

	json "github.com/goccy/go-json"
)

type JSONErr struct {
	Err string `json:"error"`
}

func (j JSONErr) Error() string {
	return j.Err
}

func NewJSONErr(err any) JSONErr {
	switch err := err.(type) {
	case string:
		return JSONErr{Err: err}
	case JSONErr:
		return JSONErr{Err: err.Err}
	case error:
		return JSONErr{Err: err.Error()}
	}
	return JSONErr{Err: fmt.Sprintf("%v", err)}
}

func JSONOut(w http.ResponseWriter, v any) {
	if err, ok := v.(error); ok {
		v = NewJSONErr(err)
	}
	json.NewEncoder(w).Encode(v)
}

func JSONOutS(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	JSONOut(w, v)
}
