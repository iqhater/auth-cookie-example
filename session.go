package main

const AuthTypeCookie = "session"
const AuthTypeJWTToken = "token"

// Session struct store key-value pair of cookie session
type Session struct {
	ID       string
	Name     string
	AuthType string
}

// NewSession function init new session struct with empty values
func NewSession() *Session {
	return &Session{
		ID:   "",
		Name: "",
	}
}
