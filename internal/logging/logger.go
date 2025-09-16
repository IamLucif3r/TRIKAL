package logging

import (
	"os"

	"github.com/charmbracelet/log"
)

type LogConfig struct {
	Level string
	JSON  bool
}

func New(cfg LogConfig) *log.Logger {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		Prefix:          "trikal",
	})

	if cfg.JSON {
		l.SetFormatter(log.JSONFormatter)
	} else {
		l.SetFormatter(log.TextFormatter)
	}

	switch cfg.Level {
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "warn":
		l.SetLevel(log.WarnLevel)
	case "error":
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}
	return l
}
