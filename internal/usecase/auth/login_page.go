package auth

import (
	"net/http"
)

func (a *AuthUsecase) SignInPage(w http.ResponseWriter, req *http.Request) error {
	return a.authRepo.SignInPage(w, req)
}
