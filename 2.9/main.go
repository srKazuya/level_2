package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	NumError = errors.New("некорректная строка, в строке только цифры")
)

func stringUnpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

	var res strings.Builder
	back := make([]rune, 0)

	for k, v := range s {
		if k == 0 && unicode.IsNumber(v) {
			return "", NumError
		}

		if !unicode.IsNumber(v) {
			back = append(back, v)
			continue
		}

		last := back[len(back)-1]
		dig, _ := strconv.Atoi(string(v))

		for i := 1; i < dig; i++ {
			res.WriteRune(last)
		}

		res.WriteRune(last)
		back = back[:0]
	}

	if len(back) > 0 {
		res.WriteString(string(back))
	}

	return res.String(), nil
}

func main() {
	tests := []string{"a4bc2d5e", "abcd", "45", ""}

	for _, t := range tests {
		res, err := stringUnpack(t)
		if err != nil {
			fmt.Printf("input=%q error: %v\n", t, err)
		} else {
			fmt.Printf("input=%q %q\n", t, res)
		}
	}
}
