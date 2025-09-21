package main

import (
	"fmt"
	"sortUtil/internal/reader"
	"sortUtil/internal/sorter"
	"sortUtil/internal/unpack"
	"sortUtil/utils/filename"
)

func main() {
	ft := reader.Reader()
	fmt.Println(ft)

	fileCounter, err := unpack.Unpack(&ft)
	if err != nil {
		fmt.Println(err)
		return
	}

	files := make([]string, fileCounter)

	for i := 0; i < fileCounter; i++ {
		filename := fmt.Sprintf("%v_%d.txt", filename.MyName(ft.File), i)
		err := sorter.SortChunk(filename, ft.Column, ft.Numeric, ft.Reverse, ft.Unique)
		if err != nil {
			fmt.Println("Ошибка сортировки:", err)
			return
		}
		files[i] = filename
	}

	outputFile := fmt.Sprintf("final_%s", ft.File)
	err = sorter.MergeFiles(files, outputFile, ft.Column, ft.Numeric, ft.Reverse, ft.Unique)
	if err != nil {
		fmt.Println("Ошибка слияния:", err)
		return
	}

	fmt.Println("Файлы успешно объединены в", outputFile)
}
