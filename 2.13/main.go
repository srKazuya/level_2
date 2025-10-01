package main

import (
	"bufio"
	"cutUtil/internal/reader"
	"fmt"
	"log"
	"os"
	"strings"
)

func ApplyCut(line string, ct reader.CutTab) string {
	if ct.Separated && !strings.Contains(line, ct.Delimiter) {
		return ""
	}

	cols := strings.Split(line, ct.Delimiter)
	var result []string

	for _, f := range ct.Fields {
		if f > 0 && f <= len(cols) {
			result = append(result, cols[f-1])
		}
	}

	for _, r := range ct.Ranges {
		start, end := r[0], r[1]
		for i := start; i <= end; i++ {
			if i > 0 && i <= len(cols) {
				result = append(result, cols[i-1])
			}
		}
	}

	return strings.Join(result, ct.Delimiter)
}

func RunCut(ct reader.CutTab) {
	file, err := os.Open(ct.File)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		out := ApplyCut(line, ct)
		if out != "" {
			fmt.Println(out)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}
}

func main() {
	ct := reader.Reader()
	RunCut(ct)
}
