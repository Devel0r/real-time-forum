package main

import (
	"fmt"
	"log/slog"

	"github.com/Pruel/real-time-forum/pkg/config"
	"github.com/Pruel/real-time-forum/pkg/logger"
)

func main() {
	// TODO: configuration
	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error(err.Error()) // log with error level// slog.Debug, Info -> os.Exit()
		// log.Fatal() -> // fmt.Println(), os.Exit(1) - process crashed
		// os.Exit(1)// os -> sys call -> exit ->
		return
	}

	// error, panic() -> defer func() { // ok
	// 	if r := recover(); r != nil {
	// 		log.Print("panic defered")
	// 	}
	// }

	// TODO: level loggining
	if _, err := logger.InitLogger(cfg); err != nil {
		slog.Error(err.Error())
		return
	}

	// TODO: App.Run, run the app logic with error handling

	// TODO: App.Stop with error handling

	fmt.Print("Hello world")
}
