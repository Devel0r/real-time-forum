package main

import (
	"log/slog"
	"os"

	"github.com/Pruel/real-time-forum/internal/app"
	"github.com/Pruel/real-time-forum/pkg/config"
	"github.com/Pruel/real-time-forum/pkg/logger"
)

func main() {
	slog.Info("Real Time Forum running...")
	// configuration
	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error(err.Error()) // log with error level// slog.Debug, Info -> os.Exit()
		// log.Fatal() -> // fmt.Println(), os.Exit(1) - process crashed
		os.Exit(1) // os -> sys call -> exit ->
		return
	}
	slog.Debug("Configuration successful initialized")

	// error, panic() -> defer func() { // ok
	// 	if r := recover(); r != nil {
	// 		log.Print("panic defered")
	// 	}
	// }

	// level loggining
	if _, err := logger.InitLogger(cfg); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Debug("Logger successful initialized")

	slog.Info("Real Time Forum started!")
	// App.Run, run the app logic with error handling
	if err := app.Run(cfg); err != nil {
		slog.Error(err.Error())
	}

	// App.Stop with error handling

}
