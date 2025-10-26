// Package telnet provides
package telnet

import (
	"errors"
	"log/slog"
	"telnet/internal/infrastructure/tcp"
)

type Service interface {
	Run(timeout int, tokens []string)
}

type service struct {
	log    *slog.Logger
	dialer tcp.Dialer
}

func NewService(log *slog.Logger, dialer tcp.Dialer) Service {
	return service{log: log, dialer: dialer}
}

var (
	errWrongSyntx = errors.New("ошибка синтаксиса команды")
)

func (s service) Run(timeout int, tokens []string) {
	op := "telnet.Service.Run"
	log := s.log
	log.With(
		slog.String("op", op),
	)

	if len(tokens) != 2 {
		log.Error("Synt error", slog.Any("error", errWrongSyntx))
	}
	host := tokens[0]
	port := tokens[len(tokens)-1]
	log.Info("Разбиение токена:", slog.Any("host", host), slog.Any("port", port), slog.Any("timeout", timeout))
	s.dialer.Connect(timeout, host, port)


}
