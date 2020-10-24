package kernel

import "syscall"

type SyscallArguments struct {
	Rdi uintptr
	Rsi uintptr
	Rdx uintptr
	R10 uintptr
	R8  uintptr
	R9  uintptr
}

type Syscall func(SyscallArguments) (uintptr, uintptr, syscall.Errno)

type SyscallTable struct {
	table map[uintptr]Syscall
}

func NewSyscallTable() *SyscallTable {
	s := SyscallTable{}
	s.table = make(map[uintptr]Syscall)
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
	return func(args SyscallArguments) (uintptr, uintptr, syscall.Errno) {
		return syscall.Syscall6(
			uintptr(id), args.Rdi,
			args.Rsi, args.Rdx,
			args.R10, args.R8, args.R9)
	}
}
