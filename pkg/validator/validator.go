package validator

import (
	"errors"
	"regexp"

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

// regexp = `^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*])[A-Za-z0-9!@#$%^&*]{8,}$` for validating password
// email =

// ValidateEmail
func ValidateEmail(email string) (valid bool) {
	const emailRegexp = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Регулярные выражение - это мощный инструмент используемый для поиска, сопаставления, и валидации строк. Рег.выр юзаются для сложных проверок например проверка emai
	regexp := regexp.MustCompile(emailRegexp)

	if ok := regexp.MatchString(email); !ok {
		return false
	}
	return true
}

func ValidatePassword(password string) (valid bool) {
	const passwordRegexp = `^[a-zA-Z\d!@#$%^&*]{8,}$`

	regexp := regexp.MustCompile(passwordRegexp)

	if ok := regexp.MatchString(password); !ok {
		return false
	}
	return true
}
