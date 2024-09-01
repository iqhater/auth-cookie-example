package main

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// create new cookie struct
func (s *Session) createCookie() *http.Cookie {

	id := uuid.NewV4()

	c := &http.Cookie{
		Name:     "session",
		Value:    id.String(),
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
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c)

	s.ID = ""
	s.Name = ""
}
