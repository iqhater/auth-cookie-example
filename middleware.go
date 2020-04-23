package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func isAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// check session if user already logged in
		c, err := req.Cookie("session")
		if err != nil {

			// check method post
			if req.Method == http.MethodPost {

				// set new cookie session id
				c = createCookie()
				sessionID = c.Value
				http.SetCookie(w, c)
			} else {
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}
		}

		if sessionID != c.Value && sessionName != c.Name {
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		next(w, req)
	}
}

func validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// check password
		err := bcrypt.CompareHashAndPassword(passwordHash, []byte(req.FormValue("password")))
		if err != nil {
			http.Redirect(w, req, "/error", http.StatusSeeOther)
			return
		}

		// check login
		login := req.FormValue("login")
		if login != LOGIN {
			http.Redirect(w, req, "/error", http.StatusSeeOther)
			return
		}
		next(w, req)
	}
}

/* func isAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
} */

func showLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("%s %s %s %s\n", req.Method, req.RemoteAddr, req.URL.Path, req.Proto)
		next(w, req)
		return
	}
}

func notFound(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		routes := []string{"/", "/login", "/logout", "/error", "/user"}

		for _, route := range routes {
			if req.URL.Path == route {
				next(w, req)
				return
			}
		}
		http.Redirect(w, req, "/404", http.StatusSeeOther)
	}
}

func createCookie() *http.Cookie {

	id := uuid.NewV4()

	return &http.Cookie{
		Name:  "session",
		Value: id.String(),
		// Secure: true, // https
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute), // 1 minute expire session removed
	}
}
