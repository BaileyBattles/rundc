package syscalls

import (
	"fmt"
	"syscall"
)

func Write(args SyscallArguments) (uintptr, uintptr, syscall.Errno) {
	buff := make([]byte, args.Rdx)
	syscall.PtracePeekData(args.Pid, args.Rsi, buff)
	fmt.Print(string(buff[:]))
	return args.Rdx, uintptr(0), syscall.Errno(uintptr(0))
}
