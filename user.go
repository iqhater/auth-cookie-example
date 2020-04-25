package main

// User security data from .env file
type User struct {
	login        string
	password     string
	passwordHash []byte
}
