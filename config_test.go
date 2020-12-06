package main

import (
	"testing"
)

func TestNewConfigNotEmptyData(t *testing.T) {

	cfg := NewConfig()

	if cfg.user == "" || cfg.password == "" || cfg.httpPort == "" || cfg.httpsPort == "" || cfg.tlsCert == "" || cfg.tlsKey == "" || string(cfg.passwordHash) == "" {
		t.Errorf("Config struct should not have an empty values: got %v", cfg)
	}
}
