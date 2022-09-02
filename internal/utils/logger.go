package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"os"
	"time"
)

func NewLogger(cfg *config.Config) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("app_name", cfg.App.Name).
		Str("host_ip", cfg.Rest.Host).
		Str("host_port", fmt.Sprint(cfg.Rest.Port)).
		Logger()

	return &logger
}
