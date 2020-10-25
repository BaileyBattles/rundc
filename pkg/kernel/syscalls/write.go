package syscalls

import (
	"fmt"
	"syscall"
)

func Write(args SyscallArguments) (uintptr, uintptr, syscall.Errno) {
	fmt.Println("you do what i say")
	return args.Rdx, uintptr(0), syscall.Errno(uintptr(1))
}
