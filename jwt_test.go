package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestCreateJWTValid(t *testing.T) {

	s := Session{
		ID:   "",
		Name: "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": s.Name,
			"exp":      time.Now().Add(time.Minute).Unix(),
		})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY") + randomSalt))

	c := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Minute), // 1 minute expire session removed
	}

	result := s.createToken()

	// check returned cookie structure has no error
	if result == nil {
		t.Errorf("createToken method returned not nil: got %v want %v", result, "filled structure with created name and value.")
	}

	// check name and value is setup correct
	if s.ID != result.Value && s.ID != "" {
		t.Errorf("createToken method wrong session value setup: got %v want %v", s.ID, c.Value)
	}

	// check name and value is setup correct
	if s.Name != result.Name && s.Name != "" {
		t.Errorf("createToken method wrong session name setup: got %v want %v", s.Name, c.Name)
	}
}

func TestDeleteJWTCorrect(t *testing.T) {

	s := Session{
		ID:       "",
		Name:     "token",
		AuthType: AuthTypeJWTToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": s.Name,
			"exp":      time.Now().Add(time.Minute).Unix(),
		})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY") + randomSalt))

	c := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now(),
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

	s.deleteToken(rr)

	// check correct header setup
	expected := ""

	if header := rr.Result().Header.Get("Set-Cookie"); header != expected {
		t.Errorf("handler returned wrong header set: got %v want %v", header, expected)
	}

	// check value is setup correct
	if s.ID != "" {
		t.Errorf("deleteToken method wrong session value setup: got %v want %v", s.ID, c.Value)
	}

	// check name is setup correct
	if s.Name != "" {
		t.Errorf("deleteToken method wrong session name setup: got %v want %v", s.Name, c.Name)
	}
}
