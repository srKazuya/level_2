package myNtp

import (
	"errors"
	"fmt"
	"os"
	"github.com/beevik/ntp"
)

var (
	ErrTime = errors.New("Ошибка получения времени ntp")
)

func MyTime() {

	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", ErrTime, err)
		os.Exit(0)
	}

	fmt.Println(time.Format("15.04.05 2006-01-02"))
}
