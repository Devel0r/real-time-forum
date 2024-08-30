package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pruel/real-time-forum/internal/controller"
	"github.com/Pruel/real-time-forum/internal/controller/router"
	"github.com/Pruel/real-time-forum/internal/controller/server"
	"github.com/Pruel/real-time-forum/pkg/cstructs"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

func Run(cfg *cstructs.Config) error {
	db, err := sqlite.InitDatabase(cfg)
	if err != nil {
		return err
	}

	if err := db.SQLite.Ping(); err != nil {
		fmt.Println("AMONG ASS")
	}

	slog.Debug("Successfuly connected to the SQLite3 database")

	ctl := controller.New(db)
	router := router.New(ctl)
	router.InitRouter()
	slog.Debug("Sucessful initialized router and created controllers")

	server, err := server.New(cfg, router)
	if err != nil {
		return err
	}
	slog.Debug("Sucessfuly create server instance")

	go func() {
		if err := server.RunServer(); err != nil {

			slog.Error(err.Error())
		}
	}()

	fmt.Printf("\nReal Time Forum started: http://%s:%s \n\tFor stoping press Ctrl + C\n", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	// TODO: remove after, Client

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sig := <-sigChan

	if err := server.StopServer(context.Background()); err != nil {
		return err
	}
	slog.Debug("Successful close http server")

	if err := db.Close(); err != nil {
		return err
	}
	slog.Debug("Successful close database connection")

	slog.Info("service stoping with signal", "sig", sig.String())

	return nil
}
