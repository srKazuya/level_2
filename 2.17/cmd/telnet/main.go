package main

import (
	"log/slog"
	"os"
	"telnet/internal/adapters/cli"
	"telnet/internal/telnet"
	"telnet/internal/infrastructure/tcp"
)

func main() {
	//===logger===
	log := setupLogger()
	log.Info("Telnet is running")
	log.Debug("logger debug mode enabled")

	dialer :=	tcp.NewDialer(log)
	service := telnet.NewService(log, dialer) // Service
	cli := cli.NewCli(log, service) //Cli


	cli.Run()

}

func setupLogger() *slog.Logger {
	var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
