package logging

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
)

const serviceName = "apollo"

// New creates a structured JSON logger configured for Apollo.
func New(output io.Writer, level string) (zerolog.Logger, error) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse log level %q: %w", level, err)
	}

	logger := zerolog.New(output).
		Level(parsedLevel).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return logger, nil
}
