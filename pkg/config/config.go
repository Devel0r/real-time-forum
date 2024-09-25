package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/Pruel/real-time-forum/pkg/cstructs"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/validator"
)

// pkg -> pkg1 -> pkg2 dataM
// pkg1 -> pkg -> pkg2 dataM

const (
	defaultServiceName = "real-time-forum"

	defaultHTTPServerHost         = "local"
	defaultHTTPServerPort         = "8080"
	defaultHTTPServerIdleTimeout  = 60 * time.Second
	defaultHTTPServerWriteTimeout = 30 * time.Second
	defaultHTTPServerReadTimeout  = 30 * time.Second
	defaultHTTPServerMaxHeaderMB  = 20 << 20

	defaultLoggerLevel     = -4
	defaultLoggerSourceKey = true
	defaultLoggerOutput    = "stdout"
	defaultLoggerHandler   = "json"
)

const (
	configDir  = "configs"
	configFile = "config.yaml"
)

// go vet (go-errors-check, go-unused, go-staticcheck, )

func InitConfig() (*cstructs.Config, error) {
	cfg := &cstructs.Config{}
	populateConfig(cfg)

	cfgPath, ok := os.LookupEnv("CONFIG_FILE_PATH")
	if !ok {
		if cfgPath == "" {
			cfgPath = strings.Join([]string{configDir, configFile}, "/")
		}
		return nil, serror.ErrInvalidConfigPath
	}

	cfg, err := parseConfigFileAndSetConfigParams(cfgPath, cfg)
	if err != nil {
		return nil, err
	}

	if err := validator.ValidateConfigParams(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfigFileAndSetConfigParams(filePath string, cfg *cstructs.Config) (*cstructs.Config, error) {
	if filePath == "" {
		return nil, serror.ErrInvalidConfigPath
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", serror.ErrFileNotExists, err)
	}

	if strings.HasSuffix(filePath, ".yaml") {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	dbFilePath, ok := os.LookupEnv("DATABASE_FILE_PATH")
	if !ok {
		return nil, serror.ErrEmptyEnv
	}

	cfg.DatabaseFilePath = dbFilePath

	return cfg, nil
}

func populateConfig(cfg *cstructs.Config) {
	cfg.ServiceName = defaultServiceName

	cfg.HTTPServer.Host = defaultHTTPServerHost
	cfg.HTTPServer.Port = defaultHTTPServerPort
	cfg.HTTPServer.IdleTimeout = defaultHTTPServerIdleTimeout
	cfg.HTTPServer.WriteTimeout = defaultHTTPServerWriteTimeout
	cfg.HTTPServer.ReadTimeout = defaultHTTPServerReadTimeout
	cfg.HTTPServer.MaxHeaderMB = defaultHTTPServerMaxHeaderMB

	cfg.Logger.Level = defaultLoggerLevel
	cfg.Logger.SourceKey = defaultLoggerSourceKey
	cfg.Logger.Output = defaultLoggerOutput
	cfg.Logger.Handler = defaultLoggerHandler
}
