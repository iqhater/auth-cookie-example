package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestIsAuthMiddleware(t *testing.T) {

	s := Session{}

	req, err := http.NewRequest(http.MethodGet, "/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	isAuthHandler := func(w http.ResponseWriter, r *http.Request) {

		c, err := req.Cookie("session")
		if err != nil {
			t.Errorf("Cookie with name 'session' not found in request: got %q", err)
		}

		// check user is already logged in
		if s.ID != c.Value || s.Name != c.Name {
			t.Errorf("Session id or session's name are equal: got %q", err)
		}
	}

	rr := httptest.NewRecorder()

	handler := s.isAuth(isAuthHandler)
	handler(rr, req)

	expected := s.ID
	cookie, _ := req.Cookie("session")

	if cookie.String() != expected {
		t.Errorf("handler returned wrong status code: got %v want %v", cookie, expected)
	}
}

func TestIsAuthMiddlewareCookieNotSet(t *testing.T) {

	s := Session{}

	req, err := http.NewRequest(http.MethodPost, "/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	isAuthHandler := func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, s.createCookie())
	}

	rr := httptest.NewRecorder()

	handler := s.isAuth(isAuthHandler)
	handler(rr, req)

	name := "session"
	cookie, _ := req.Cookie(name)
	val := rr.Result().Cookies()[0].Value

	if s.ID != val || s.Name != name {
		t.Errorf("handler returned wrong status code: got %v want %v", cookie, s.ID)
	}

	if rr.Result().StatusCode != http.StatusSeeOther {
		t.Error("Bad redirect status!")
	}
}

func TestIsAuthMiddlewareBadCookie(t *testing.T) {

	s := Session{}

	rr := httptest.NewRecorder()
	http.SetCookie(rr, s.createCookie())

	req := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}

	isAuthHandler := func(w http.ResponseWriter, r *http.Request) {

		name := "session"
		fakeName := "fake_session"
		cookie, _ := r.Cookie(name)
		val := cookie.Value + "fake"

		if s.Name == fakeName || s.ID == val {
			t.Errorf("handler returned wrong status code: got %v want %v", val, s.ID)
		}
	}

	handler := s.isAuth(isAuthHandler)
	handler(rr, req)

	if rr.Result().StatusCode == http.StatusSeeOther {
		t.Error("Bad redirect status!")
	}
}

func TestValidateMiddleware(t *testing.T) {

	pass := "somepass"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	u := User{
		login:        "testuser",
		password:     pass,
		passwordHash: passwordHash,
	}

	buf := new(bytes.Buffer)
	params := url.Values{}
	params.Set("login", u.login)
	params.Set("password", pass)
	buf.WriteString(params.Encode())

	req, err := http.NewRequest(http.MethodPost, "/login", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	validateHandler := func(w http.ResponseWriter, r *http.Request) {

		err := bcrypt.CompareHashAndPassword(u.passwordHash, []byte(req.FormValue("password")))

		if req.Method != http.MethodPost || err != nil || req.FormValue("login") != u.login {
			t.Errorf("Wrong method %v. Login or password incorrect! %v %v", req.Method, err, req.FormValue("login"))
		}
	}

	rr := httptest.NewRecorder()

	handler := u.validate(validateHandler)
	handler(rr, req)
}

func TestValidateMiddlewareInvalid(t *testing.T) {

	pass := "somepass"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	u := User{
		login:        "testuser",
		password:     pass,
		passwordHash: passwordHash,
	}

	buf := new(bytes.Buffer)
	params := url.Values{}
	params.Set("login", "wrong")
	params.Set("password", "bad")
	buf.WriteString(params.Encode())

	req, err := http.NewRequest(http.MethodPost, "/login", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	validateHandler := func(w http.ResponseWriter, r *http.Request) {

		err := bcrypt.CompareHashAndPassword(u.passwordHash, []byte(req.FormValue("password")))

		if req.Method == http.MethodPost || err == nil || req.FormValue("login") == u.login {
			t.Errorf("Wrong method %v. Error shuold be non nil. Login or password incorrect! %v %v", req.Method, err, req.FormValue("login"))
		}
	}

	rr := httptest.NewRecorder()

	handler := u.validate(validateHandler)
	handler(rr, req)

	if rr.Result().StatusCode != http.StatusSeeOther {
		t.Error("Bad redirect status!")
	}
}

func TestShowLogMiddleware(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	showLogHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s %s\n", req.Method, req.RemoteAddr, req.URL.Path, req.Proto)
	}

	rr := httptest.NewRecorder()

	handler := showLog(http.HandlerFunc(showLogHandler))
	handler(rr, req)

	if buf.Len() == 0 {
		t.Error("Empty log output!")
	}
}

func TestNotFoundMiddlewareRedirect(t *testing.T) {

	r := Routes{
		routes: []string{"/", "/one", "/two", "/three"},
	}

	req, err := http.NewRequest(http.MethodGet, "/wqer", nil)
	if err != nil {
		t.Fatal(err)
	}

	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {

		for _, route := range r.routes {
			if req.URL.Path == route {
				return
			}
		}
	}

	rr := httptest.NewRecorder()

	handler := r.notFound(notFoundHandler)
	handler(rr, req)

	if rr.Result().StatusCode != http.StatusSeeOther {
		t.Error("Wrong redirect status code!")
	}
}

func TestNotFoundMiddleware(t *testing.T) {

	r := Routes{
		routes: []string{"/", "/one", "/two", "/three"},
	}

	req, err := http.NewRequest(http.MethodGet, "/two", nil)
	if err != nil {
		t.Fatal(err)
	}

	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {

		for _, route := range r.routes {
			if req.URL.Path == route {
				return
			}
		}
	}

	rr := httptest.NewRecorder()

	handler := r.notFound(notFoundHandler)
	handler(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Wrong status code!")
	}
}

func TestSecureFilesNotAuthorizedMiddleware(t *testing.T) {

	// s := &Session{}

	r := Routes{
		// Session: s,
		files: []string{"user.html", "gopher_wizard.png"},
	}

	req, err := http.NewRequest(http.MethodGet, "/public/img/gopher_wizard.png", nil)
	if err != nil {
		t.Fatal(err)
	}

	secureFilesHandler := func(w http.ResponseWriter, req *http.Request) {

		for _, file := range r.files {
			if path.Base(req.URL.Path) == file {

				_, err := req.Cookie("session")
				if err == nil {
					t.Errorf("Cookie with name 'session' should not exist in request: got %q", err)
				}

				// check user is already logged in
				/* if r.ID != c.Value || r.Name != c.Name {
					t.Errorf("Session id or session's name are equal: got %q", err)
				} */
				return
			}
		}
	}

	rr := httptest.NewRecorder()

	handler := r.secureFiles(http.HandlerFunc(secureFilesHandler))
	// secureFilesHandler(rr, req)
	handler.ServeHTTP(rr, req)

	// expected := r.ID
	// cookie, _ := req.Cookie("session")

	/* if cookie.String() != expected {
		t.Errorf("handler returned wrong status code: got %v want %v", cookie, expected)
	} */

	if rr.Result().StatusCode != http.StatusSeeOther {
		t.Error("Wrong redirect status code!")
	}
}

func TestSecureFilesNotInListMiddleware(t *testing.T) {

	r := Routes{
		files: []string{"gopher_wizard.png"},
	}

	req, err := http.NewRequest(http.MethodGet, "/public/img/fake_wizard.png", nil)
	if err != nil {
		t.Fatal(err)
	}

	secureFilesHandler := func(w http.ResponseWriter, req *http.Request) {

		for _, file := range r.files {
			if path.Base(req.URL.Path) == file {

				_, err := req.Cookie("session")
				if err == nil {
					t.Errorf("Cookie with name 'session' should not exist in request: got %q", err)
				}
				return
			}
		}
	}

	rr := httptest.NewRecorder()

	handler := r.secureFiles(http.HandlerFunc(secureFilesHandler))
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Wrong redirect status code!", rr.Result().StatusCode)
	}
}

func TestSecureFilesAuthorizedMiddleware(t *testing.T) {

	s := &Session{}

	r := Routes{
		Session: s,
		files:   []string{"user.html", "gopher_wizard.png"},
	}

	req, err := http.NewRequest(http.MethodGet, "/public/img/gopher_wizard.png", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.SetCookie(rr, r.createCookie())

	req.Header.Set("Cookie", rr.HeaderMap["Set-Cookie"][0])

	secureFilesHandler := func(w http.ResponseWriter, req *http.Request) {

		for _, file := range r.files {
			if path.Base(req.URL.Path) == file {

				c, err := req.Cookie("session")
				if err != nil {
					t.Errorf("Cookie with name 'session' not exist in request: got %q", err)
				}

				//TODO: cover test case
				if r.ID != c.Value || r.Name != c.Name {
					t.Errorf("Session id or session's name are equal: got %q", err)
				}
				return
			}
		}
	}

	handler := r.secureFiles(http.HandlerFunc(secureFilesHandler))
	handler.ServeHTTP(rr, req)

	expected := r.ID
	cookie := rr.Result().Cookies()[0].Value

	if cookie != expected {
		t.Errorf("handler returned wrong cookie: got %v want %v", cookie, expected)
	}

	/* if rr.Result().StatusCode != http.StatusSeeOther {
		t.Error("Wrong redirect status code!")
	} */
}

func TestSecureHeadersMiddleware(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	secureHeadersHandler := func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
	}

	rr := httptest.NewRecorder()

	handler := secureHeaders(http.HandlerFunc(secureHeadersHandler))
	handler.ServeHTTP(rr, req)

	// check correct header setup
	expected := "1; mode=block"

	if header := rr.Result().Header.Get("X-XSS-Protection"); header != expected {
		t.Errorf("handler returned wrong header set: got %v want %v", header, expected)
	}

	expected2 := "deny"

	if header := rr.Result().Header.Get("X-Frame-Options"); header != expected2 {
		t.Errorf("handler returned wrong header set: got %v want %v", header, expected2)
	}
}
