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

type CutTab struct {
	Fields     []int  // -f
	Ranges     [][2]int // диапазоны полей
	Delimiter  string // -d
	Separated  bool   // -s
	File       string 
}

var (
	ErrEmptyInLine = errors.New("Строка не должна быть пустой:\n")
	ErrUnknownFlg  = errors.New("Неизвестный флаг: ")
)

func parseFields(fieldsArg string) ([]int, [][2]int) {
	var fields []int
	var ranges [][2]int

	parts := strings.Split(fieldsArg, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			b := strings.Split(p, "-")
			if len(b) == 2 {
				start, err1 := strconv.Atoi(b[0])
				end, err2 := strconv.Atoi(b[1])
				if err1 == nil && err2 == nil && start <= end {
					ranges = append(ranges, [2]int{start, end})
				}
			}
		} else {
			if n, err := strconv.Atoi(p); err == nil {
				fields = append(fields, n)
			}
		}
	}
	return fields, ranges
}

func Reader() CutTab {
	fmt.Println("Введите флаги и файл для cut (например: cut -f 1,3-5 -d , -s file.txt):")
	in := bufio.NewReader(os.Stdin)

	inLine, _ := in.ReadString('\n')
	inLine = strings.TrimSpace(inLine)

	args := strings.Fields(inLine)
	if len(args) == 0 {
		log.Fatal(ErrEmptyInLine)
	}

	var ct CutTab
	ct.Delimiter = "\t" 

	if args[0] == "cut" {
		for i := 1; i < len(args); i++ {
			elem := args[i]

			if strings.HasPrefix(elem, "-") {
				flag := elem[1:]
				switch flag {
				case "f":
					if i+1 < len(args) {
						fields, ranges := parseFields(args[i+1])
						ct.Fields = fields
						ct.Ranges = ranges
						i++
					}
				case "d":
					if i+1 < len(args) {
						ct.Delimiter = args[i+1]
						i++
					}
				case "s":
					ct.Separated = true
				default:
					log.Fatal(ErrUnknownFlg, flag)
				}
			} else {
				if i == len(args)-1 {
					ct.File = elem
				}
			}
		}
	}
	return ct
}
