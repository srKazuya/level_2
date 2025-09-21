package reader

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	ErrFlagPars   = errors.New("Строка должна содержать флаги (начинаться с '-'):\n")
	ErrFileOpen   = errors.New("Ошибка открытия файла:\n")
	ErrUnknownFlg = errors.New("Неизвестный флаг: ")
	ErrKFlag      = errors.New("Флаг -k требует параметр — номер колонки (например: -k 2):\n")
)

type FTab struct {
	Kol     bool // k
	Numeric bool // n
	Reverse bool // r
	Unique  bool // u
	Column  int  // k N
	File    string
}


func Reader() FTab {
	fmt.Println("Введите флаги и файл для сортировки (например: -k 2 -nr file.txt):")
	in := bufio.NewReader(os.Stdin)

	inLine, _ := in.ReadString('\n')
	inLine = strings.TrimSpace(inLine)

	args := strings.Fields(inLine)
	if len(args) == 0 {
		log.Fatal(ErrFlagPars)
	}

	var ft FTab
	for i := 0; i < len(args); i++ {
		a := args[i]

		if strings.HasPrefix(a, "-") {
			flags := a[1:]
			j := 0
			for j < len(flags) {
				switch flags[j] {
				case 'k':
					numStr := ""
					if j+1 < len(flags) { 
						numStr = flags[j+1:]
						j = len(flags) 
					} else if i+1 < len(args) {
						numStr = args[i+1]
						i++ 
					} else {
						log.Fatal(ErrKFlag)
					}
					val, err := strconv.Atoi(numStr)
					if err != nil {
						log.Fatal(ErrKFlag)
					}
					ft.Column = val
					ft.Kol = true
				case 'n':
					ft.Numeric = true
				case 'r':
					ft.Reverse = true
				case 'u':
					ft.Unique = true
				default:
					log.Fatal(ErrUnknownFlg, string(flags[j]))
				}
				j++
			}
		} else {
			ft.File = a
		}
	}

	if ft.File == "" {
		log.Fatal("Не указан файл для сортировки")
	}

	return ft
}
