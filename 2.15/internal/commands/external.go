package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"minishell/internal/commands/cmdErrs"
	"os"
	"os/exec"
	"syscall"
)

func ExternalCmd(cLine []string) string {
	if len(cLine) == 0 {
		fmt.Println("Строка должна содеражить команду command <args>")
		return ""
	}
	cmd := cLine[0]
	extCmd := exec.Command(cmd, cLine[1:]...)

	var out bytes.Buffer
	var errBuf bytes.Buffer
	extCmd.Stdout = &out
	extCmd.Stderr = &errBuf
	extCmd.Stdin = os.Stdin

	err := extCmd.Run()
	if err != nil {
		if IsInterrupted(err) {
			fmt.Printf("Команда %v прервана пользователем\n", cmd)
			return ""
		}
		fmt.Println(cmdErrs.ErrExtCommand, err)
		if errBuf.Len() > 0 {
			fmt.Println(errBuf.String())
		}
		return ""
	}
	return out.String()
}

func IsInterrupted(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}

	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return false
	}
	return status.Signaled() && status.Signal() == syscall.SIGINT
}
// TODO:LEARN
func ExternalCmdWithInput(args []string, input io.Reader) string {
    cmd := exec.Command(args[0], args[1:]...)
    if input != nil {
        cmd.Stdin = input
    }
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error:", err)
    }
    return out.String()
}