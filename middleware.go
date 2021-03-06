package main

import (
	"log"
	"net/http"
	"path"

	"golang.org/x/crypto/bcrypt"
)

// isAuth middleware handler check cookie session
// and redirect to the main page if cookie is not setup or invalid credentials
func (s *Session) isAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// check session if user already logged in
		c, err := req.Cookie("session")
		if err != nil {

			// check method post
			if req.Method == http.MethodPost {

				// set new cookie session id
				http.SetCookie(w, s.createCookie())

				// redirect to user page
				http.Redirect(w, req, "/user", http.StatusSeeOther)
				return
			}
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		// check user is already logged in
		if s.ID != c.Value || s.Name != c.Name {
			s.deleteCookie(w)
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		next(w, req)
	}
}

// validate middleware handler check user login and password.
// If user login and password are wrong then redirect to the error page
func (u *User) validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// parse login
		login := req.FormValue("login")

		// check password
		err := bcrypt.CompareHashAndPassword(u.passwordHash, []byte(req.FormValue("password")))

		// check method post, check login and check password
		if req.Method != http.MethodPost || err != nil || login != u.login {
			http.Redirect(w, req, "/error", http.StatusSeeOther)
			return
		}
		next(w, req)
	}
}

// showLog middleware handler shows network data log info
func showLog(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s %s\n", req.Method, req.RemoteAddr, req.URL.Path, req.Proto)
		next.ServeHTTP(w, req)
	})
}

// notFound middleware handler check existing routes.
// If route was not found in list then redirect to the 404 page
func (r *Routes) notFound(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		for _, route := range r.routes {
			if req.URL.Path == route {
				next(w, req)
				return
			}
		}
		http.Redirect(w, req, "/404", http.StatusSeeOther)
	}
}

// secureFiles middleware handler protect your public folder files from unauthorized users
// by check session and corresponding files and redirect to 404 page
func (r *Routes) secureFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		for _, file := range r.files {

			if path.Base(req.URL.Path) == file {

				// check session if user already logged in
				c, err := req.Cookie("session")
				if err != nil {
					http.Redirect(w, req, "/404", http.StatusSeeOther)
					return
				}

				// check user is already logged in
				if r.ID != c.Value || r.Name != c.Name {
					r.deleteCookie(w)
					http.Redirect(w, req, "/404", http.StatusSeeOther)
					return
				}
				next.ServeHTTP(w, req)
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}

// secureHeaders middleware handler setup secure headers
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, req)
	})
}
