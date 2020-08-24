package main

import (
	"net/http"
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

	s.ID = c.Value
	s.Name = c.Name

	result := s.createCookie()

	// check returned cookie structure has no error
	if result == nil {
		t.Errorf("createCookie method returned not nil: got %v want %v", result, "filled structure with created name and value.")
	}

	// check name nad value is setup correct
	if s.ID != result.Value && s.ID != "" {
		t.Errorf("createCookie method wrong session value setup: got %v want %v", s.ID, c.Value)
	}

	// check name nad value is setup correct
	if s.Name != result.Name && s.Name != "" {
		t.Errorf("createCookie method wrong session name setup: got %v want %v", s.Name, c.Name)
	}
}

func TestDeleteCookieCorrect(t *testing.T) {

	//TODO:
}
