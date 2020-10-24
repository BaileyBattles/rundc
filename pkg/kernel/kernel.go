package kernel

import "rundc/pkg/abi"

type Kernel struct {
	Table *SyscallTable
}

func (k *Kernel) HandleSyscall(process *Process, id uintptr, args SyscallArguments) {
	if syscall := k.Table.GetSyscall(id); syscall != nil {
		syscall(args)
		abi.SetSyscall(process.cmd.Process.Pid, 64)
	}
	err := process.Ptrace(24)
	process.WaitForStatus()
	if err != nil {
		panic("Failed")
	}

}
