package builtins

import (
	"fmt"
	"minishell/internal/commands/cmdErrs"
	"os"
)

func Cd(cLine []string) {
	if len(cLine) < 2 {
		fmt.Println("Использование: cd <path>")
	}

	path := cLine[1]
	err := os.Chdir(path)
	if err != nil {
		fmt.Println(cmdErrs.ErrChangeDir, err)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(cmdErrs.ErrChangeDir, err)
	}
	fmt.Println("Текущая директория: ", dir)
}
