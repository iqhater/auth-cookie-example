package main

// User security data from .env file
type User struct {
	login        string
	password     string
	passwordHash []byte
}

// NewUser function takes a Config structure and return a new User struct
func NewUser(config *Config) *User {
	return &User{
		login:        config.user,
		password:     config.password,
		passwordHash: config.passwordHash,
	}
}
