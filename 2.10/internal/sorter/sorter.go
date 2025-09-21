package sorter

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getKey(line string, column int) string {
	parts := strings.Split(line, "\t")
	if column > 0 && column <= len(parts) {
		return parts[column-1]
	}
	return line
}

func SortChunk(filename string, column int, numeric, reverse, unique bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	sort.Slice(lines, func(i, j int) bool {
		ai := getKey(lines[i], column)
		aj := getKey(lines[j], column)

		if numeric {
			af, errA := strconv.ParseFloat(ai, 64)
			bf, errB := strconv.ParseFloat(aj, 64)

			if errA != nil {
				af = 0
			}
			if errB != nil {
				bf = 0
			}

			if reverse {
				return af > bf
			}
			return af < bf
		}

		if reverse {
			return ai > aj
		}
		return ai < aj
	})

	if unique {
		uniq := make([]string, 0, len(lines))
		seen := make(map[string]struct{})
		for _, l := range lines {
			if _, ok := seen[l]; !ok {
				seen[l] = struct{}{}
				uniq = append(uniq, l)
			}
		}
		lines = uniq
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	for _, l := range lines {
		if _, err := w.WriteString(l + "\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

func MergeFiles(files []string, output string, column int, numeric, reverse, unique bool) error {
	type fileLine struct {
		line string
		key  string
		f    *os.File
		r    *bufio.Scanner
	}

	instances := make([]*fileLine, 0, len(files))
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		fl := &fileLine{f: f, r: scanner}
		if scanner.Scan() {
			fl.line = scanner.Text()
			fl.key = getKey(fl.line, column)
			instances = append(instances, fl)
		} else {
			f.Close()
		}
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bufio.NewWriter(out)

	seen := make(map[string]struct{})

	for len(instances) > 0 {
		selected := 0
		for i := 1; i < len(instances); i++ {
			if compareKeys(instances[i].key, instances[selected].key, numeric, reverse) {
				selected = i
			}
		}

		if !unique || (unique && seen[instances[selected].line] == struct{}{}) {
			w.WriteString(instances[selected].line + "\n")
			if unique {
				seen[instances[selected].line] = struct{}{}
			}
		}

		if instances[selected].r.Scan() {
			instances[selected].line = instances[selected].r.Text()
			instances[selected].key = getKey(instances[selected].line, column)
		} else {
			instances[selected].f.Close()
			instances = append(instances[:selected], instances[selected+1:]...)
		}
	}

	return w.Flush()
}

func compareKeys(a, b string, numeric, reverse bool) bool {
	if numeric {
		af, _ := strconv.ParseFloat(a, 64)
		bf, _ := strconv.ParseFloat(b, 64)
		if reverse {
			return af > bf
		}
		return af < bf
	}
	if reverse {
		return a > b
	}
	return a < b
}
