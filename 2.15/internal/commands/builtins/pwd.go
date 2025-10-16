package builtins
 
import (
	"fmt"
	"minishell/internal/commands/cmdErrs"
	"os"
)

func Pwd(cLine []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err, cmdErrs.ErrGetwd)
	}
	fmt.Println(dir)
}
