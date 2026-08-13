package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/config"
)

func New(cfg *config.LogConfig, timeFormat string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf(`parse log level %q: %w`, cfg.Level, err)
	}

	var output io.Writer = os.Stdout

	switch strings.ToLower(cfg.Format) {
	case "console":
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		}
	case "json":
		// JSON is default format
	default:
		return zerolog.Logger{}, fmt.Errorf("unsupported log format %q: expected console or json", cfg.Format)
	}

	zerolog.TimeFieldFormat = timeFormat
	return zerolog.New(output).Level(level).With().Timestamp().Logger(), nil
}
