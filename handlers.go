package main

import (
	"net/http"
)

func login(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/login.html")
}

func logout(w http.ResponseWriter, req *http.Request) {

	// delete cookie
	c := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func admin(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/user.html")
}

func unAuth(w http.ResponseWriter, req *http.Request) {
	// w.WriteHeader(http.StatusUnauthorized)
	http.ServeFile(w, req, "./public/error.html")
}
