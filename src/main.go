package main

import (
	"fmt"

	"os"
	"os/exec"
	"syscall"
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
	cmd := exec.Command("./a.out")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned with err: %v\n\n\n", err.Error())
	}

	for {
		if _, _, errno := syscall.Syscall6(syscall.SYS_PTRACE, PTRACE_SYSCALL, uintptr(cmd.Process.Pid), 0, 0, 0, 0); errno != 0 {
			fmt.Printf("Ptrace failed with errno = %s\n", errno.Error())
		}

		var status syscall.WaitStatus
		_, err = syscall.Wait4(cmd.Process.Pid, &status, 0, nil)

		if err != nil {
			fmt.Printf("error in wait 4 %s\n", err.Error())
		}

		var regs syscall.PtraceRegs
		err = syscall.PtraceGetRegs(cmd.Process.Pid, &regs)
		if err != nil {
			fmt.Printf("Error getting regs: %s\n", err.Error())
		}

		if status.Stopped() && status.StopSignal() == syscall.SIGSEGV {
			fmt.Println("Stopped")
			fmt.Println(status.StopSignal())
			syscall.Syscall6(syscall.SYS_PTRACE, syscall.PTRACE_CONT, uintptr(cmd.Process.Pid), 0, 0, 0, 0)
		}

		// if shouldCall {
		// 	fmt.Printf("Runing syscall %s\n", syscallName(regs.Orig_rax))
		// 	fmt.Printf("Regs used = %+v\n", regs)
		// 	fmt.Println(regs.Orig_rax, regs.Rdi, regs.Rsi, regs.Rdx, regs.R10, regs.R8, regs.R9)
		// 	r1, r2, errno := syscall.Syscall6(uintptr(regs.Orig_rax), uintptr(regs.Rdi), uintptr(regs.Rsi), uintptr(regs.Rdx), uintptr(regs.R10), uintptr(regs.R8), uintptr(regs.R9))
		// 	syscall.PtraceGetRegs(cmd.Process.Pid, &regs)
		// 	if errno != 0 {
		// 		fmt.Printf("Forwarding syscall error = %d\n\n\n\n\n", errno)
		// 	} else {
		// 		fmt.Print("\n\n\n\n\n")
		// 	}
		// 	fmt.Sprintf("%d%d", r1, r2)

		// }
		// shouldCall = !shouldCall

		if err != nil {
			os.Exit(2)
		}
	}
}

// err = syscall.PtraceSyscall(cmd.Process.Pid, 0)
// if err != nil {
// 	fmt.Println("error with ptrace")
// }
