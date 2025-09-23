package main

import (
	"fmt"
	"sort"
	"strings"

)

type res map[string][]string

func keyGen(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	return string(runes)
}

func anagram(input []string) res {
	seen := make(map[string]bool)      // слово - bool
	group := make(map[string][]string) // ключ - слова
	first := make(map[string]string)   // ключ - первое слово

	for _, w := range input {
		w = strings.ToLower(w)
		if seen[w] {
			continue
		}
		seen[w] = true

		key := keyGen(w)

		if _, ok := first[key]; !ok {
			first[key] = w
		}
		group[key] = append(group[key], w)
	}
	result := make(res)
	for key, words := range group {
		if len(words) < 2 {
			continue
		}
		sort.Strings(words)
		result[first[key]] = words
	}
	return result
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	res := anagram(input)

	fmt.Println(res)
}
