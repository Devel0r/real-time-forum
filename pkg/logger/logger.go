package logger

import (
	"log/slog"
	"os"

	"github.com/Pruel/real-time-forum/pkg/cstructs"
	"github.com/Pruel/real-time-forum/pkg/validator"
)

func InitLogger(cfg *cstructs.Config) (logger *slog.Logger, err error) {
	if err := validator.ValidateConfigParams(cfg); err != nil {
		return nil, err
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.Level(cfg.Logger.Level),
			AddSource: cfg.Logger.SourceKey,
		}).WithAttrs([]slog.Attr{
			slog.String("service_name", cfg.ServiceName)}))

	slog.SetDefault(logger)

	return logger, nil
}
