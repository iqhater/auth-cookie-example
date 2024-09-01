package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)

// login handler serve login.html page
func login(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/login.html")
}

// logout handler remove cookie and redirect to main page
func (s *Session) logout(w http.ResponseWriter, req *http.Request) {
	s.deleteCookie(w)
	s.deleteToken(w)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// admin handler set cache-control header and serve user.html page
func admin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")

	tpl, err := template.ParseFiles("./public/user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	authType := strings.Split(req.Header.Get("Cookie"), "=")

	renderValues := make(map[string]string)

	if authType[0] == "session" {
		renderValues["Title"] = "Cookie"
		renderValues["Image"] = "/public/img/cookie_logo.svg"
		renderValues["Time"] = "1 minute"
	}

	if authType[0] == "token" {
		renderValues["Title"] = "Token"
		renderValues["Image"] = "/public/img/jwt_logo.svg"
		renderValues["Time"] = "5 minutes"
	}

	err = tpl.ExecuteTemplate(w, "user", renderValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// unAuth handler serve error.html page
func unAuth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	http.ServeFile(w, req, "./public/error.html")
}

// pageNotFound handler serve 404.html page
func pageNotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, req, "./public/404.html")
}

// redirectToHTTPS handler get hostname
// and redirect to https schema with 301 status code
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {

	host := strings.Split(req.Host, ":")[0] + ":"
	httpsPort := os.Getenv("HTTPS_PORT")

	http.Redirect(w, req, "https://"+host+httpsPort, http.StatusMovedPermanently)
}
