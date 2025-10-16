package  builtins
  
import (
	"fmt"
	"minishell/internal/commands/cmdErrs"
	"os"
	"strconv"
)

func Kill(cLine []string) {
	if len(cLine) < 1 {
		return
	}
	pid, err := strconv.Atoi(cLine[0])
	if err != nil {
		fmt.Println(err)
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(cmdErrs.ErrExtCommand, err)
	}

	err = p.Kill()
	if err != nil {
		fmt.Println(cmdErrs.ErrExtCommand, err)
	}
	fmt.Println("Процесс убит: ", pid)
}
