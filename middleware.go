package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/tabwriter"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Session) setCookie(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// check method post
		if req.Method == http.MethodPost {

			// set new cookie session id
			if s.AuthType == AuthTypeCookie {
				http.SetCookie(w, s.createCookie())
			}

			// set new jwt token
			if s.AuthType == AuthTypeJWTToken {
				http.SetCookie(w, s.createToken())
			}

			fmt.Println("User page redirect...")
			// successful redirect to user page
			http.Redirect(w, req, "/user", http.StatusSeeOther)
			return
		} else {
			// w.WriteHeader(http.StatusMethodNotAllowed)
			http.Redirect(w, req, "/", http.StatusSeeOther)
		}
		next(w, req)
	}
}

// isAuth middleware handler check cookie session
// and redirect to the main page if cookie is not setup or invalid credentials
func (s *Session) isAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		log.Println("AUTH TYPE: ", s.AuthType)
		// check session if user already logged in
		c, err := req.Cookie(s.AuthType)
		if err != nil {

			log.Println("Cookie Error: ", err)
			if err == http.ErrNoCookie {

				// redirect to login page or 401 status on same page?
				// w.WriteHeader(http.StatusUnauthorized)
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}
		}

		// check user is already logged in for cookie session
		if s.AuthType == AuthTypeCookie {
			if s.ID != c.Value || s.Name != c.Name {
				s.deleteCookie(w)
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}
		}

		// check user is already logged in for jwt
		if s.AuthType == AuthTypeJWTToken {
			if verifyToken(c.Value) != nil {
				s.deleteToken(w)
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}
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

// setAuthType middleware handler setup auth type
func (s *Session) setAuthType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// logout every time as check login & password validation
		// otherwise session & token can live in cookies together at the same time.
		// USE ONLY with /login route!
		s.deleteCookie(w)
		s.deleteToken(w)

		// parse query method type "session" or "token"
		authType := req.FormValue("authType")

		// set query to Session struct
		if authType == AuthTypeCookie {
			s.AuthType = AuthTypeCookie
		}

		if authType == AuthTypeJWTToken {
			s.AuthType = AuthTypeJWTToken
		}
		fmt.Println("Auth Type: ", s.AuthType)
		next(w, req)
	}
}

// showLog middleware handler shows network data log info
func showLog(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()

		sr := NewStatusHTTP(w)
		next.ServeHTTP(sr, req)

		statusCode := sr.StatusCode

		tw := tabwriter.NewWriter(os.Stdout, 28, 4, 1, ' ', tabwriter.Debug)
		fmt.Fprintf(tw, "%v\t [%d: %s]\t %v\t %s\t %s\t %s\n", t.Format("02.01.2006 15:04:05"), statusCode, http.StatusText(statusCode), time.Since(t), req.RemoteAddr, req.Method, req.URL.String())
		tw.Flush()
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
				c, err := req.Cookie(r.AuthType)
				if err != nil {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				// check user is already logged in for cookie session
				if r.AuthType == AuthTypeCookie {
					if r.ID != c.Value || r.Name != c.Name {
						r.deleteCookie(w)
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
				}

				// check user is already logged in for jwt
				if r.AuthType == AuthTypeJWTToken {
					if verifyToken(c.Value) != nil {
						r.deleteToken(w)
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
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
