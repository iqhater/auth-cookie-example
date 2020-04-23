package main

import (
	"net/http"
)

func login(w http.ResponseWriter, req *http.Request) {

	// deleteCookie(w)
	http.ServeFile(w, req, "./public/login.html")
}

func logout(w http.ResponseWriter, req *http.Request) {

	deleteCookie(w)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func admin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	http.ServeFile(w, req, "./public/user.html")
}

func unAuth(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/error.html")
}

func pageNotFound(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/404.html")
}
