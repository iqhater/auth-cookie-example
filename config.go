package main

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Config struct contains init environments
// variables and generated passwordHash
type Config struct {
	user         string
	password     string
	httpPort     string
	httpsPort    string
	tlsCert      string
	tlsKey       string
	forcedTLS    bool
	passwordHash []byte
}

// NewConfig function returns inited server configuration
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

	forcedTLSString, exists := os.LookupEnv("FORCED_TLS")
	if !exists {
		log.Println("Env variable FORCED_TLS does not exist!")
	}

	forcedTLS, err := strconv.ParseBool(forcedTLSString)
	if err != nil {
		log.Println(err)

		// set default to false
		forcedTLS = false
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
		forcedTLS,
		passwordHash,
	}
}
