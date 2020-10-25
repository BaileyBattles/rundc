package syscalls

type SyscallArguments struct {
	Rdi uintptr
	Rsi uintptr
	Rdx uintptr
	R10 uintptr
	R8  uintptr
	R9  uintptr
}
