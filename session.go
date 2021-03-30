package main

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Session struct store key-value pair of cookie session
type Session struct {
	ID   string
	Name string
}

// create new cookie struct
func (s *Session) createCookie() *http.Cookie {

	id := uuid.NewV4()

	c := &http.Cookie{
		Name:     "session",
		Value:    id.String(),
		Path:     "/user",
		Secure:   true, // https
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Minute), // 1 minute expire session removed
	}

	s.ID = c.Value
	s.Name = c.Name
	return c
}

// deleteCookie method remove cookie and reset user credentials
func (s *Session) deleteCookie(w http.ResponseWriter) {

	c := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c)

	s.ID = ""
	s.Name = ""
}

// NewSession function init new session struct with empty values
func NewSession() *Session {
	return &Session{
		ID:   "",
		Name: "",
	}
}
