// Package cli
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"telnet/internal/telnet"

	"github.com/google/shlex"
	"github.com/spf13/pflag"
)

type CLI interface {
	Run()
	Manage(lineTokens []string) error
}

type cli struct {
	log     *slog.Logger
	service telnet.Service
}

func NewCli(log *slog.Logger, service telnet.Service) CLI {
	return &cli{
		log:     log,
		service: service,
	}
}

var (
	errReadInString = errors.New("ошибка чтения строки")
	errSplitInLine  = errors.New("ошибка токенизации строки")
	errFlagParsing  = errors.New("ошибка парсинга флагов")
	errEmptyLine    = errors.New("пустая строка")
)

func (c *cli) Run() {
	//===logger===
	const op = "controllers.CLI.RUN"
	log := c.log
	log = log.With(
		slog.String("op", op),
	)

	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Введите: flags <host> [port]")
		fmt.Print("telnet>")

		//===Read input===
		inLine, err := readInput(in)
		if errors.Is(err, errReadInString) {
			log.Error("ошибка ввода", slog.Any("error", err))
			continue
		}

		//===Split line in tokens===
		lineTokens, err := shlex.Split(inLine)
		if err != nil {
			log.Error(errSplitInLine.Error(), slog.Any("error", err))
			continue
		}

		// ===Manage Command===
		if err = c.Manage(lineTokens); err != nil {
			log.Error("Manage error", slog.Any("error", err))
		}
	}

}

func (c *cli) Manage(tokens []string) error {
	fs := pflag.NewFlagSet("telnet", pflag.ContinueOnError)
	timeout := fs.IntP("timeout", "t", 0, "set timeout")
	if err := fs.Parse(tokens); err != nil {
		return fmt.Errorf("%w: %v", errFlagParsing, err)
	}

	out := strings.Join(fs.Args(), " ")
	if out == "" {
		return errEmptyLine
	}

	tokens = strings.Fields(out)

	flagsUsed := fs.Changed("timeout")

	if !flagsUsed {
		fmt.Println("Нет флага")
		c.service.Run(*timeout, tokens)
	}
	if fs.Changed("timeout") {
		fmt.Printf("флаг %v\n", timeout)
		c.service.Run(*timeout, tokens)
	}
	return nil
}

func readInput(in *bufio.Reader) (string, error) {
	inLine, err := in.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("%w: %w", errReadInString, err)
	}
	inLine = strings.TrimSpace(inLine)
	return inLine, nil
}
