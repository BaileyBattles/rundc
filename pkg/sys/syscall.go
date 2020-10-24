package sys

import (
	"fmt"

	sec "github.com/seccomp/libseccomp-golang"
)

func PrintSyscallName(id uint64) {
	name, _ := sec.ScmpSyscall(id).GetName()
	fmt.Printf("Called syscall: %s\n", name)
}
