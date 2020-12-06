package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	user         string
	password     string
	httpPort     string
	httpsPort    string
	tlsCert      string
	tlsKey       string
	passwordHash []byte
}

func NewConfig() *Config {

	user, exists := os.LookupEnv("LOGIN")
	if !exists {
		log.Println("Env variable LOGIN does not exist!")
	}

	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		log.Println("Env variable PASSWORD does not exist!")
	}

	httpPort, exists := os.LookupEnv("HTTP_PORT")
	if !exists {
		log.Println("Env variable HTTP_PORT does not exist!")
	}

	httpsPort, exists := os.LookupEnv("HTTPS_PORT")
	if !exists {
		log.Println("Env variable HTTPS_PORT does not exist!")
	}

	tlsCert, exists := os.LookupEnv("TLS_CERT_PATH")
	if !exists {
		log.Println("Env variable TLS_CERT_PATH does not exist!")
	}

	tlsKey, exists := os.LookupEnv("TLS_KEY_PATH")
	if !exists {
		log.Println("Env variable TLS_KEY_PATH does not exist!")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		user,
		password,
		httpPort,
		httpsPort,
		tlsCert,
		tlsKey,
		passwordHash,
	}
}
