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
	secretKey    string
	httpPort     string
	httpsPort    string
	tlsCert      string
	tlsKey       string
	forcedTLS    bool
	passwordHash []byte
}

// NewConfig function returns inited server configuration
func NewConfig() *Config {

	user := getEnv("LOGIN")
	password := getEnv("PASSWORD")
	secretKey := getEnv("SECRET_KEY")
	httpPort := getEnv("HTTP_PORT")
	httpsPort := getEnv("HTTPS_PORT")
	tlsCert := getEnv("TLS_CERT_PATH")
	tlsKey := getEnv("TLS_KEY_PATH")
	forcedTLSString := getEnv("FORCED_TLS")

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
		secretKey,
		httpPort,
		httpsPort,
		tlsCert,
		tlsKey,
		forcedTLS,
		passwordHash,
	}
}

// getEnv wrapper function to get a value from environment
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Env variable %s does not exist!\n", key)
		return ""
	}
	return value
}
