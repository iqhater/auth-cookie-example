package auth

import "net/http"

func (s *AuthSession) SignInPage(w http.ResponseWriter, req *http.Request) error {
	http.ServeFile(w, req, "./public/login.html")
	return nil
}
