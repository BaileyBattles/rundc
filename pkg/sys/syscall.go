package sys

import (
	"fmt"

	sec "github.com/seccomp/libseccomp-golang"
)

const (
	PTRACE_SYSEMU  = 31
	PTRACE_SYSCALL = 24
)

func PrintSyscallName(id uint64) {
	name, _ := sec.ScmpSyscall(id).GetName()
	fmt.Printf("Called syscall: %s\n", name)
}
