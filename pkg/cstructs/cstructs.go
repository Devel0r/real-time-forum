package cstructs

import "time"



type (
	HTTPServer struct {
		Host         string        `json:"host" yaml:"host"`
		Port         string        `json:"port" yaml:"port"`
		IdleTimeout  time.Duration `json:"idle_time" yaml:"idle_time"`
		WriteTimeout time.Duration `json:"write_time" yaml:"write_time"`
		ReadTimeout  time.Duration `json:"read_time" yaml:"read_time"`
		MaxHeaderMB  int           `json:"max_header_mb" yaml:"max_header_mb"`
	}
	Logger struct {
		Level     int    `json:"level" yaml:"level"`
		SourceKey bool   `json:"source_key" yaml:"source_key"`
		Output    string `json:"output" yaml:"output"`
		Handler   string `json:"handler" yaml:"handler"`
	}

	Config struct {
		ServiceName string `json:"service_name" yaml:"service_name"`
		HTTPServer  `json:"http_server" yaml:"http_server"`
		Logger      `json:"logger" yaml:"logger"`
	}
)