package main

import (
	"testing"
)

func TestNewSessionEmptyValues(t *testing.T) {

	s := NewSession()

	if s.ID != "" || s.Name != "" {
		t.Errorf("Session struct should be init with empty values: got %v", s)
	}
}
