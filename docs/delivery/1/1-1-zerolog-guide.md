# 1-1 zerolog Integration Guide

Date: 2026-02-18

## Package

- `github.com/rs/zerolog`

## Documentation

- https://github.com/rs/zerolog

## API Usage Notes

1. Create a logger with `zerolog.New(io.Writer).With().Timestamp().Logger()`.
2. Emit JSON logs with `logger.Info().Str("key", "value").Msg("message")`.
3. Set log level globally with `zerolog.SetGlobalLevel(level)`.
4. Use `zerolog.ParseLevel("info")` when converting string config to level.

## Apollo Usage Example

```go
package logging

import (
    "os"

    "github.com/rs/zerolog"
)

func NewLogger(level string) (zerolog.Logger, error) {
    parsedLevel, err := zerolog.ParseLevel(level)
    if err != nil {
        return zerolog.Logger{}, err
    }

    logger := zerolog.New(os.Stdout).
        Level(parsedLevel).
        With().
        Timestamp().
        Str("service", "apollo").
        Logger()

    return logger, nil
}
```

## Key Patterns

- Build one base logger at startup and pass child loggers as needed.
- Keep message text concise and push details into structured fields.
- Always terminate chained events with `.Msg(...)` or `.Send()`.
