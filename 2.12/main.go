package main

import (
	"bufio"
	"errors"
	"fmt"
	"grepUtil/internal/reader"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	ErrFileOpen = errors.New("Ошибка открытия файла:\n")
)

func main() {
	gt := reader.Reader()
	fmt.Printf("%+v\n", gt)

	if err := search(&gt); err != nil {
		log.Fatal(err)
	}

}

func search(gt *reader.GrepTab) error {
	file, err := os.Open(gt.File)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFileOpen, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if gt.C.N > 0 {
		gt.A.N, gt.B.N = gt.C.N, gt.C.N
	}


	var re *regexp.Regexp
	pattern := gt.SeParam
	if !gt.F { 
		flags := ""
		if gt.I {
			flags = "(?i)"
		}
		re = regexp.MustCompile(flags + pattern)
	} else if gt.I {
		pattern = strings.ToLower(pattern)
	}

	matchLines := make([]bool, len(lines))
	matchesCount := 0
	for i, line := range lines {
		matched := false
		if gt.F {
			text := line
			if gt.I {
				text = strings.ToLower(text)
			}
			matched = strings.Contains(text, pattern)
		} else {
			matched = re.MatchString(line)
		}
		if gt.V {
			matched = !matched
		}
		if matched {
			matchLines[i] = true
			matchesCount++
		}
	}

	if gt.SmC {
		fmt.Println(matchesCount)
		return nil
	}

	printed := make([]bool, len(lines))
	for i := range lines {
		if matchLines[i] {
			start := i - gt.B.N
			if start < 0 {
				start = 0
			}
			end := i + gt.A.N
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for j := start; j <= end; j++ {
				if !printed[j] {
					if gt.N {
						fmt.Printf("%d:", j+1)
					}
					fmt.Println(lines[j])
					printed[j] = true
				}
			}
		}
	}

	return nil
}