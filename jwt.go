package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var randomSalt = generateRandomToken(4)

// createToken creates a new token
func (s *Session) createToken() *http.Cookie {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": s.Name,
			"exp":      time.Now().Add(5 * time.Minute).Unix(),
		})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY") + randomSalt))
	if err != nil {
		return nil
	}

	c := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Secure:   true, // https
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(5 * time.Minute), // 5 minutes expire session removed
	}

	s.ID = c.Value
	s.Name = c.Name

	return c
}

// deleteToken method remove jwt token and reset user credentials
func (s *Session) deleteToken(w http.ResponseWriter) {

	// immediately clear the token cookie
	c := &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now(),
	}

	http.SetCookie(w, c)

	s.ID = ""
	s.Name = ""
}

func verifyToken(rawToken string) error {

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY") + randomSalt), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}

func generateRandomToken(length uint) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
