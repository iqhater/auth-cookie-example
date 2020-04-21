package main

import "net/http"

var sessionName string
var sessionID string

// delete cookie
func deleteCookie(w http.ResponseWriter) {

	c := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	sessionName = ""
	sessionID = ""
}
