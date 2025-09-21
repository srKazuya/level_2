package unpack

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sortUtil/internal/reader"
)

var (
	ErrOpenFile = errors.New("Ошибка при открытии файла")
)

const chunkLim = 10 * 1024

var store map[string]struct{}

func Unpack(r *reader.FTab) (int, error) {
	inFile, err := os.Open(r.File)
	if err != nil {
		return 0, fmt.Errorf("%w : %v", ErrOpenFile, err)
	}
	defer inFile.Close()

	f := ""
	fmt.Println(f)

	store = make(map[string]struct{})
	in := bufio.NewReader(inFile)

	counter := 0
	totalSize := 0
	for {
		line, err := in.ReadString('\n')
		if len(line) > 0 {
			store[line] = struct{}{}
			totalSize += len(line) + 16
			fmt.Printf("Элементов в store: %d, totalSize=%d\n", len(store), totalSize)

			if totalSize >= chunkLim {
				if err := flushStore(r.File, store, counter, totalSize); err != nil {
					return 0, err
				}
				counter++
				store = make(map[string]struct{})
				totalSize = 0
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
	}

	if len(store) > 0 {
		if err := flushStore(r.File, store, counter, totalSize); err != nil {
			return 0, err
		}
	}

	return counter, nil
}

func flushStore(fileName string, store map[string]struct{}, counter int, size int) error {
	ext := filepath.Ext(fileName)
	name := fileName[:len(fileName)-len(ext)]
	filename := fmt.Sprintf("%v_%d.txt", name, counter)
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%w : %v", ErrOpenFile, err)
	}
	defer outFile.Close()

	out := bufio.NewWriter(outFile)
	for st := range store {
		if _, err := out.WriteString(st); err != nil {
			return err
		}
	}
	out.Flush()

	fmt.Println("Сохранил", filename, "размером ~", size, "байт")
	return nil
}
