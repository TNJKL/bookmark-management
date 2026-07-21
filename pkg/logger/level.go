package logger

import "github.com/rs/zerolog"

// SetLogLevel sets the global logging level for zerolog based on the provided level string.
func SetLogLevel(levelStr string) {
	level := zerolog.NoLevel
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.NoLevel
	}
	zerolog.SetGlobalLevel(level)
}
