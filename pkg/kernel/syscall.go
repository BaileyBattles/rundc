package kernel

import (
	"rundc/pkg/kernel/syscalls"
	"syscall"
)

type Syscall func(syscalls.SyscallArguments) (uintptr, uintptr, syscall.Errno)

type SyscallTable struct {
	table map[uintptr]Syscall
}

func NewSyscallTable() *SyscallTable {
	s := SyscallTable{}
	s.table = map[uintptr]Syscall{
		uintptr(1): syscalls.Write,
	}
	return &s
}

func (s *SyscallTable) Register(id uintptr, f Syscall) {
	s.table[id] = f
}

func (s *SyscallTable) GetSyscall(id uintptr) Syscall {
	if val, ok := s.table[id]; ok {
		return val
	}
	return nil
}

func basicSyscall(id uintptr) Syscall {
	return func(args syscalls.SyscallArguments) (uintptr, uintptr, syscall.Errno) {
		return syscall.Syscall6(
			uintptr(id), args.Rdi,
			args.Rsi, args.Rdx,
			args.R10, args.R8, args.R9)
	}
}
