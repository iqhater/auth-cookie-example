package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestCreateCookieValid(t *testing.T) {

	s := Session{
		ID:   "",
		Name: "",
	}

	id := uuid.NewV4()

	c := http.Cookie{
		Name:     "session",
		Value:    id.String(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Minute), // 1 minute expire session removed
	}

	result := s.createCookie()

	// check returned cookie structure has no error
	if result == nil {
		t.Errorf("createCookie method returned not nil: got %v want %v", result, "filled structure with created name and value.")
	}

	// check name and value is setup correct
	if s.ID != result.Value && s.ID != "" {
		t.Errorf("createCookie method wrong session value setup: got %v want %v", s.ID, c.Value)
	}

	// check name and value is setup correct
	if s.Name != result.Name && s.Name != "" {
		t.Errorf("createCookie method wrong session name setup: got %v want %v", s.Name, c.Name)
	}
}

func TestDeleteCookieCorrect(t *testing.T) {

	s := Session{
		ID:       uuid.NewV4().String(),
		Name:     "session",
		AuthType: AuthTypeCookie,
	}

	c := http.Cookie{
		Name:     "session",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	s.ID = c.Value
	s.Name = c.Name

	req, err := http.NewRequest(http.MethodGet, "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(http.NotFound)

	handler(rr, req)

	s.deleteCookie(rr)

	// check correct header setup
	expected := ""

	if header := rr.Result().Header.Get("Set-Cookie"); header != expected {
		t.Errorf("handler returned wrong header set: got %v want %v", header, expected)
	}

	// check value is setup correct
	if s.ID != "" {
		t.Errorf("deleteCookie method wrong session value setup: got %v want %v", s.ID, c.Value)
	}

	// check name is setup correct
	if s.Name != "" {
		t.Errorf("deleteCookie method wrong session name setup: got %v want %v", s.Name, c.Name)
	}
}
