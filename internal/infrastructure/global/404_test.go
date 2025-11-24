package global

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPageNotFoundOK(t *testing.T) {

	// Arrange
	globalRepo := NewGlobalRepoImpl()

	req, err := http.NewRequest(http.MethodGet, "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Act
	result := globalRepo.PageNotFound(rr, req)

	// Assert
	if result != nil {
		t.Errorf("SignIn method returned not nil: got %v want %v", result, nil)
	}
}
