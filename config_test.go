package main

import (
	"os"
	"testing"
)

func TestNewConfigNotEmptyData(t *testing.T) {

	cfg := NewConfig()

	if cfg.user == "" || cfg.password == "" || cfg.httpPort == "" || cfg.httpsPort == "" || cfg.tlsCert == "" || cfg.tlsKey == "" || string(cfg.passwordHash) == "" {
		t.Errorf("Config struct should not have an empty values: got %v", cfg)
	}
}

func TestNewConfigEmptyData(t *testing.T) {

	envs := []string{"LOGIN", "PASSWORD", "HTTP_PORT", "HTTPS_PORT", "TLS_CERT_PATH", "TLS_KEY_PATH", "FORCED_TLS"}
	envsBuffer := make(map[string]string)

	// clear environments variables
	for _, env := range envs {

		// save to buffer for later restore envs
		envsBuffer[env] = os.Getenv(env)

		// clear env
		os.Setenv(env, "")
	}

	cfg := NewConfig()

	if cfg.user != "" || cfg.password != "" || cfg.httpPort != "" || cfg.httpsPort != "" || cfg.tlsCert != "" || cfg.tlsKey != "" || cfg.forcedTLS {
		t.Errorf("Config struct should be an empty values: got %v", cfg)
	}

	// restore envs for another tests
	// loads values from .env into the system
	for k, v := range envsBuffer {
		os.Setenv(k, v)
	}
}

func TestNewConfigEnvsNotExist(t *testing.T) {

	envs := []string{"LOGIN", "PASSWORD", "HTTP_PORT", "HTTPS_PORT", "TLS_CERT_PATH", "TLS_KEY_PATH", "FORCED_TLS"}
	envsBuffer := make(map[string]string)

	// clear environments variables
	for _, env := range envs {

		// save to buffer for later restore envs
		envsBuffer[env] = os.Getenv(env)

		// clear env
		os.Setenv(env, "")
	}

	// flush all environments
	// os.Clearenv()

	for _, env := range envs {

		// clear env
		os.Setenv(env, "")

		value := getEnv(env)
		if value != "" {
			t.Errorf("Env variable %s should not exist!", env)
		}
	}

	// restore envs for another tests
	// loads values from .env into the system
	for k, v := range envsBuffer {
		os.Setenv(k, v)
	}
}

func TestGetEnvExist(t *testing.T) {

	envs := []string{"LOGIN", "PASSWORD", "HTTP_PORT", "HTTPS_PORT", "TLS_CERT_PATH", "TLS_KEY_PATH", "FORCED_TLS"}

	for _, env := range envs {

		value := getEnv(env)
		if value == "" {
			t.Errorf("Env variable %s does not exist!", env)
		}
	}
}

func TestGetEnvNotExist(t *testing.T) {

	envs := []string{"FAKE_LOGIN", "FAKE_PASSWORD", "FAKE_PORT"}

	for _, env := range envs {

		// clear env
		os.Setenv(env, "")

		value := getEnv(env)
		if value != "" {
			t.Errorf("Env variable %s should not exist!", env)
		}
	}
}
