package main

import (
	"fmt"
	"rundc/pkg/kernel"

	"os"
	"syscall"
)

const (
	PTRACE_SYSEMU  = 31
	PTRACE_SYSCALL = 24
)

func main() {
	// p := kernel.CreatePtraceProcess("./test/test_syscallreturn.o", []string{})
	p := kernel.CreatePtraceProcess("/bin/bash", []string{})

	p.Start()
	err := p.Wait()
	if err.Error() != "stop signal: trace/breakpoint trap" {
		fmt.Printf("Wait returned with err: %v\n\n\n", err.Error())
	}

	for {
		p.Ptrace(PTRACE_SYSCALL)

		status, err := p.WaitForStatus()
		if status.Exited() {
			fmt.Println("Child process has exited")
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("Error waiting for status %s\n", err.Error())
		}

		if status.Stopped() && status.StopSignal() == syscall.SIGTRAP {

			err = p.HandleSyscall()
			if err != nil {
				fmt.Printf("Error handling syscall: %s\n", err.Error())
			}

		}
		if status.Stopped() && status.StopSignal() == syscall.SIGSEGV {
			panic("Received a signal")
		}
	}
}
