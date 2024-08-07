package validator

import (
	"errors"

	"github.com/Pruel/real-time-forum/pkg/cstructs"
)

func ValidateConfigParams(cfg *cstructs.Config) error {

	if cfg.ServiceName == "" {
		return errors.New("error, service name is empty")
	}

	if err := validateHTTPServerCfgParams(cfg); err != nil {
		return err
	}

	if err := validateLoggerCfgParams(cfg); err != nil {
		return err
	}

	return nil
}

func validateLoggerCfgParams(cfg *cstructs.Config) error {
	debugLevel := -4
	errorLevel := 8

	if cfg.Logger.Level < debugLevel || cfg.Logger.Level > errorLevel {
		return errors.New("error, logger level incorrect")
	}
	if cfg.Logger.Handler == "" {
		return errors.New("error, logger handler is empty")
	}

	// slog.TextHandler -> "error, "10.10.2024 10:10:10 GT+0: some Daniil log"
	// slog.JSONHandler -> "{ "error", 10.10.2024 10:10:10 GT+0: "some Daniil log"}

	if cfg.Logger.SourceKey == false {
		return errors.New("error, source key false")
	}
	if cfg.Logger.Output == "" {
		return errors.New("error, output empty")
	}

	return nil
}

func validateHTTPServerCfgParams(cfg *cstructs.Config) error {
	if cfg.HTTPServer.Host == "" {
		return errors.New("error, host is empty")
	}

	if cfg.HTTPServer.Port == "" {
		return errors.New("error, port is empty")
	}
	if cfg.HTTPServer.IdleTimeout <= 0 {
		return errors.New("error, HTTP server idle timeout must be positive")
	}

	if cfg.HTTPServer.WriteTimeout <= 0 {
		return errors.New("error, HTTP server write timeout must be positive")
	}

	if cfg.HTTPServer.ReadTimeout <= 0 {
		return errors.New("error, HTTP server read timeout must be positive")
	}

	if cfg.HTTPServer.MaxHeaderMB <= 0 {
		return errors.New("error, HTTP server max header size must be positive")
	}

	return nil
}
