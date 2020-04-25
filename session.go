package main

import "net/http"

// Session struct store key-value pair of cookie session
type Session struct {
	ID   string
	Name string
}

// delete cookie
func (s *Session) deleteCookie(w http.ResponseWriter) {

	c := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	s.ID = ""
	s.Name = ""
}
