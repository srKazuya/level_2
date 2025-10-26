// Package cli provides command-line parsing and handlers for the wget tool.
package cli

import (
	"errors"
	"fmt"
	"strings"
	"wget/internal/downloader"

	"github.com/spf13/pflag"
)

var (
	errFlagParsing = errors.New("ошибка парсинга флагов")
)

type Handler struct {
	service downloader.Service
}

func NewHandler(service downloader.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleCommand(tokens []string) error {
	fs := pflag.NewFlagSet("wget", pflag.ContinueOnError)
	r := fs.BoolP("recursive", "r", false, "recursive downloading")
	l := fs.IntP("level", "l", 0, "level of recursion")
	p := fs.BoolP("page-requisites", "p", false, "download <img>, <css>, <js>")

	if err := fs.Parse(tokens); err != nil {
		return fmt.Errorf("%w: %v", errFlagParsing, err)
	}

	out := strings.Join(fs.Args(), " ")
	if out == "" {
		return fmt.Errorf("не указан URL")
	}

	flagsUsed := fs.Changed("recursive") || fs.Changed("level") || fs.Changed("page-requisites")

	// no flag proc
	if !flagsUsed {
		_, err := h.service.GetPage(out)
		if err != nil {
			return err
		}
		return nil
	}
	//flag processing
	if *r {
		fmt.Println("recursive: true")
	}
	if fs.Changed("level") {
		fmt.Printf("level: %d\n", *l)
	}
	if *p {
		file, err := h.service.GetPage(out)
		if err != nil {
			return err
		}
		h.service.ParsePage(&file)
		fmt.Println("page-requisites: true")
	}

	return nil
}
