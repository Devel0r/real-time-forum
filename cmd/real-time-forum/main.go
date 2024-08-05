package main

import (
	"fmt"
	"log/slog"

	"github.com/Pruel/real-time-forum/pkg/config"
	"github.com/Pruel/real-time-forum/pkg/logger"
)

func main() {
	// TODO: configuration
	cfg, err := config.InitConfig(); 
    if err != nil {
        slog.Error(err.Error())
        return
    }

    // TODO: level loggining
    if _, err := logger.InitLogger(cfg); err != nil {
        slog.Error(err.Error())
    } 

	// TODO: App.Run, run the app logic with error handling
    

	// TODO: App.Stop with error handling

	fmt.Print("Hello world")
}
