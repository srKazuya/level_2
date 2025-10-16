package builtins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

var errUnknownFlag = errors.New("ошибка парсинга флагов")

func Echo(cLine []string) string {
	fs := pflag.NewFlagSet("echo", pflag.ContinueOnError)
	e := fs.BoolP("enable-escapes", "e", false, "включить поддержку вывода Escape последовательностей")

	err := fs.Parse(cLine[1:])
	if err != nil {
		fmt.Printf("%v\n %v", errUnknownFlag, err)
	}

	out := strings.Join(fs.Args(), " ")

	if *e {
		replacer := strings.NewReplacer(
			`\\`, `\`,
			`\n`, "\n",
			`\t`, "\t",
			`\r`, "\r",
			`\b`, "\b",
			`\a`, "\a",
			`\f`, "\f",
			`\v`, "\v",
		)
		out = replacer.Replace(out)
	}
	return out
}
