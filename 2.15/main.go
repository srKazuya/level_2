package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"minishell/internal/commands"
	"minishell/internal/utils"
	"os"
	"strings"

	"github.com/google/shlex"
)

var (
	errParseLine = errors.New("Ошибка парсинга строки:\n")
	errReadLine  = errors.New("Ошибка чтения команды:\n")
)

type CommandSet map[string]struct{}
type PrefixSet map[string]struct{}

func main() {

	in := bufio.NewReader(os.Stdin)
	utils.SigListener()

	for {
		fmt.Print("miniShell> ")
		//read
		inLine, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err, errReadLine)
		}
		inLine = strings.TrimSpace(inLine)

		lineTokens, err := shlex.Split(inLine)
		if err != nil {
			log.Fatalf("%v %v", errParseLine, err)
		}

		segments := make([][]string, 0)
		cLine := make([]string, 0, len(lineTokens))
		//parse
		for _, tVal := range lineTokens {
			if _, ok := utils.PrefixSet[tVal]; ok {
				segments = append(segments, cLine)
				cLine = nil
				continue
			}
			cLine = append(cLine, tVal)
		}
		segments = append(segments, cLine)

	
		if len(segments) > 1 {
			
			fmt.Println("Pipe segments:", len(segments))

			
			var prevOut io.Reader

			for i, cLine := range segments {
				var cmdOut bytes.Buffer

			
				if _, ok := utils.CommandSet[cLine[0]]; ok {
				
					builtinsOut := utils.CmdManage(cLine)

					if prevOut != nil {
						inputBytes, _ := io.ReadAll(prevOut)
						builtinsOut = utils.CmdManage(append(cLine, string(inputBytes)))
					}

					cmdOut.WriteString(builtinsOut)
				} else {
					
					if prevOut != nil {
						
						externalOut := commands.ExternalCmdWithInput(cLine, prevOut)
						cmdOut.WriteString(externalOut)
					} else {
						externalOut := commands.ExternalCmd(cLine)
						cmdOut.WriteString(externalOut)
					}
				}
				
				prevOut = bytes.NewBuffer(cmdOut.Bytes())

				if i == len(segments)-1 {
					io.Copy(os.Stdout, prevOut)
				}
			}

		} else {
			cLine := segments[0]
			if _, ok := utils.CommandSet[cLine[0]]; ok {
				builtinsOut := utils.CmdManage(cLine)
				fmt.Println(builtinsOut)
			} else {
				cmd := commands.ExternalCmd(cLine)
				fmt.Println(cmd)

			}
		}

	}
}
