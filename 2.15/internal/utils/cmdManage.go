package utils

import (
	"minishell/internal/commands/builtins"
)

var (
	CommandSet = map[string]struct{}{
		"echo": {},
		"cd":   {},
		"pwd":  {},
		"kill": {},
		"ps":   {},
	}
	PrefixSet = map[string]struct{}{
		"|": {},
	}
)

var cmdOutput string

func CmdManage(cLine []string) string {
	command := cLine[0]
	if _, ok := CommandSet[command]; ok {
		switch command {
		case "echo":
			cmdOutput = builtins.Echo(cLine)
		case "cd":
			builtins.Cd(cLine)
		case "ps":
			builtins.Ps()
		case "kill":
			builtins.Kill(cLine)
		}
		return cmdOutput
	}
	return ""
}
