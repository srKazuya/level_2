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

type CompoundFLag struct {
	InTable bool
	N       int
}

type GrepTab struct {
	//Flags
	A       CompoundFLag // - A N
	B       CompoundFLag // -B N
	C       CompoundFLag //-C N
	SmC     bool         // -c
	I       bool         // -i
	V       bool         // -v
	F       bool         // -F
	N       bool         //-n
	SeParam string
	File    string
}

var (
	ErrEmptyInLine = errors.New("Строка не должна быть пустой:\n")
	ErrUnknownFlg = errors.New("Неизвестный флаг: ")
)

func Reader() GrepTab {
	fmt.Println("Введите флаги и файл для поиска (например: grep [опции] 'шаблон_поиска' [файл(ы)]):")
	in := bufio.NewReader(os.Stdin)

	inLine, _ := in.ReadString('\n')
	inLine = strings.TrimSpace(inLine)

	args := strings.Fields(inLine)
	if len(args) == 0 {
		log.Fatal(ErrEmptyInLine)
	}

	var gt GrepTab
	if args[0] == "grep" {
		for i := 1; i < len(args); i++ {
			elem := args[i]

			if strings.HasPrefix(elem, "-") {
				flag := elem[1:]
				switch flag {
				case "A":
					gt.A.InTable = true
					nextArg := args[i+1]
					if flagN, err := strconv.Atoi(nextArg); err == nil {
						gt.A.N = flagN
					}
				case "B":
					gt.B.InTable = true
					nextArg := args[i+1]
					if flagN, err := strconv.Atoi(nextArg); err == nil {
						gt.B.N = flagN
					}
				case "C":
					gt.C.InTable = true
					nextArg := args[i+1]
					if flagN, err := strconv.Atoi(nextArg); err == nil {
						gt.C.N = flagN
					}
				case "c":
					gt.SmC = true
				case "i":
					gt.I = true
				case "v":
					gt.V = true
				case "F":
					gt.F = true
				case "n":
					gt.N = true
				default:
					log.Fatal(ErrUnknownFlg, flag)
				}
			}

			if i == len(args)-1 {
				gt.File = elem
			}
			if strings.HasPrefix(elem, "'") {
				gt.SeParam = elem[1 : len(elem)-1]
			}
		}
	}
	return gt
}
