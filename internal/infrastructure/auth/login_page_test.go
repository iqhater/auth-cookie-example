package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignInOK(t *testing.T) {

	// Arrange
	authRepo := NewAuthRepo()

	req, err := http.NewRequest(http.MethodPost, "/sign-in", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Act
	result := authRepo.SignInPage(rr, req)

	// Assert
	if result != nil {
		t.Errorf("SignIn method returned not nil: got %v want %v", result, nil)
	}
}
