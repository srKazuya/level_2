package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/shlex"

	"wget/internal/adapters/cli"
	"wget/internal/downloader"
)

var (
	errInLineRead = errors.New("ошибка чтения входа")
	errParseLine  = errors.New("ошибка парсинга строки")
)

func main() {
	// --- wire dependencies ---
	storage := &downloader.FileStorage{}          // реализация Storage
	service := downloader.NewService(storage)     // бизнес-логика
	handler := cli.NewHandler(service)            // CLI слой

	in := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("Введите флаги и ссылку (<flag> <https://...>)")
		fmt.Print("wget-> ")
		inLine, err := in.ReadString('\n')
		if err != nil {
			log.Fatalf("%v: %v", errInLineRead, err)
		}

		inLine = strings.TrimSpace(inLine)
		if inLine == "exit" {
			fmt.Println("Выход...")
			break
		}

		lineTokens, err := shlex.Split(inLine)
		if err != nil {
			log.Fatalf("%v: %v", errParseLine, err)
		}

		if err := handler.HandleCommand(lineTokens); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		}
	}
}
