package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandlerOK(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(login)

	handler.ServeHTTP(rr, req)

	// check the status code is 200.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}

func TestLogoutHandlerRedirectStatus(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	s := Session{
		ID:   "",
		Name: "",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.logout)

	handler.ServeHTTP(rr, req)

	// check the status code is 303.
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}

func TestAdminHandlerOK(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	// w.Header().Set("Cache-Control", "no-cache, private, max-age=0")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(admin)

	handler.ServeHTTP(rr, req)

	// check the status code is 200.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check correct header setup
	expected := "no-cache, private, max-age=0"

	if header := rr.Result().Header.Get("Cache-Control"); header != expected {
		t.Errorf("handler returned wrong header set: got %v want %v", header, expected)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}

func TestUnAuthHandlerOK(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(unAuth)

	handler.ServeHTTP(rr, req)

	// check the status code is 200.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}

func TestPageNotFoundHandlerOK(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pageNotFound)

	handler.ServeHTTP(rr, req)

	// check the status code is 200.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}

func TestRedirectToHTTPSHandlerRedirectStatus(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToHTTPS)

	handler.ServeHTTP(rr, req)

	// check the status code is 301.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}

	// check the response body is not empty.
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned unexpected body length: got %v want %v", rr.Body.Len(), "> 0")
	}
}
