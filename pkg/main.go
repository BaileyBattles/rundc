package main

import (
	"fmt"

	"os"
	"syscall"

	"rundc/pkg/process"
	"rundc/pkg/sys"
)

//handle syscall :platform/ptr

const (
	PTRACE_SYSEMU  = 31
	PTRACE_SYSCALL = 24
)

var syscallMap map[uint64]string = map[uint64]string{
	0:  "Read",
	1:  "Write",
	2:  "Open",
	3:  "Close",
	9:  "Mmap",
	12: "Brk",
	21: "Access",
}

func syscallName(num uint64) string {
	val, ok := syscallMap[num]
	if !ok {
		return fmt.Sprintf("%d", num)
	} else {
		return val
	}
}

func main() {
	p := process.CreatePtraceProcess("echo", []string{"here"})

	p.Start()
	err := p.Wait()
	if err != nil {
		fmt.Printf("Wait returned with err: %v\n\n\n", err.Error())
	}

	for {
		p.Ptrace(PTRACE_SYSCALL)

		status, err := p.WaitForStatus()
		if err != nil {
			fmt.Printf("Error waiting for status %s\n", err.Error())
		}

		regs, err := p.GetRegs()
		if err != nil {
			fmt.Printf("error getting regs: %s\n", err.Error())
		}
		sys.PrintSyscallName(regs.Orig_rax)

		if status.Stopped() && status.StopSignal() == syscall.SIGSEGV {
			fmt.Println("Stopped")
			fmt.Println(status.StopSignal())
			err = p.Ptrace(syscall.PTRACE_CONT)
		}
		if err != nil {
			os.Exit(2)
		}
	}
}
