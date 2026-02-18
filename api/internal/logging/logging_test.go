package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLoggerWithValidLevel(t *testing.T) {
	var output bytes.Buffer

	logger, err := New(&output, "info")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	logger.Info().Str("component", "test").Msg("logger configured")

	logged := output.String()
	if !strings.Contains(logged, "\"service\":\"apollo\"") {
		t.Fatalf("expected service field in log output, got %q", logged)
	}

	if !strings.Contains(logged, "\"component\":\"test\"") {
		t.Fatalf("expected component field in log output, got %q", logged)
	}
}

func TestNewLoggerWithInvalidLevel(t *testing.T) {
	_, err := New(&bytes.Buffer{}, "invalid-level")
	if err == nil {
		t.Fatalf("expected invalid log level to return an error")
	}
}
