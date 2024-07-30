package config

import "time"

const (
	defaultServiceName = "real-time-forum"

	defaultHTTPServerHost = "local"
	defaultHTTPServerPort = "8080"
	defaultHTTPServerIdleTimeout = 60 * time.Second 
	defaultHTTPServerWriteTimeout = 30 * time.Second
	defaultHTTPServerReadTimeout = 30 * time.Second
	defaultHTTPServerMaxHeaderMB = 20 << 20  

	defaultLoggerLevel = -4
	defaultLoggerSourceKey = true
	defaultLoggerOutput = "stdout" 
	defaultLoggerHandler = "json"
)

func InitConfig() {

}

func PopulateConfig() {

}