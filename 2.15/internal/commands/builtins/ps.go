package builtins

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
)

func Ps() {
	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%-8s %-8s %-8s %-s\n", "PID", "PPID", "CPU%", "COMMAND")

	for _, p := range processes {
		pid := p.Pid

		name, _ := p.Name()
		ppid, _ := p.Ppid()
		cpuPercent, _ := p.CPUPercent()

		fmt.Printf("%-8d %-8d %-8.2f %-s\n", pid, ppid, cpuPercent, name)
	}
}
