package main

import (
	"fmt"

	"os"
	"syscall"

	"rundc/pkg/process"
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
	counter := 0
	p := process.CreatePtraceProcess("../test/test_syscallreturn.o", []string{})

	p.Start()
	err := p.Wait()
	if err.Error() != "stop signal: trace/breakpoint trap" {
		fmt.Printf("Wait returned with err: %v\n\n\n", err.Error())
	}

	for {
		if counter == 44 {
			p.Ptrace(PTRACE_SYSEMU)

		} else {
			p.Ptrace(PTRACE_SYSCALL)
		}

		status, err := p.WaitForStatus()
		if status.Exited() {
			fmt.Println("Child process has exited")
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("Error waiting for status %s\n", err.Error())
		}

		if status.Stopped() && status.StopSignal() == syscall.SIGTRAP {

			err = p.HandleSyscall(counter == 44)
			if err != nil {
				fmt.Printf("Error handling syscall: %s\n", err.Error())
			}
			counter += 1

		}
		if status.Stopped() && status.StopSignal() == syscall.SIGSEGV {
			os.Exit(2)
			fmt.Println("Stopped")
			fmt.Println(status.StopSignal())
			s, _ := p.GetSignalInfo()
			fmt.Println(s)
			err = p.Ptrace(syscall.PTRACE_CONT)
		}
	}
}
