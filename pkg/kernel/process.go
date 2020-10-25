package kernel

import (
	"os/exec"
	"syscall"

	"rundc/pkg/abi"
	"rundc/pkg/kernel/syscalls"
)

type Process struct {
	kernel *Kernel
	cmd    *exec.Cmd
}

func (this *Process) Start() error {
	return this.cmd.Start()
}

func (this *Process) Wait() error {
	return this.cmd.Wait()
}

func (this *Process) Ptrace(flags int) error {
	if _, _, errno := syscall.Syscall6(syscall.SYS_PTRACE, uintptr(flags), uintptr(this.cmd.Process.Pid), 0, 0, 0, 0); errno != 0 {
		return error(errno)
	}
	return nil
}

func (this *Process) HandleSyscall() error {
	regs, err := this.getRegs()
	if err != nil {
		return err
	}
	//sys.PrintSyscallName(regs.Orig_rax)

	if err = this.kernel.HandleSyscall(this, uintptr(regs.Orig_rax), syscalls.SyscallArguments{
		Rdi: uintptr(regs.Rdi),
		Rsi: uintptr(regs.Rsi),
		Rdx: uintptr(regs.Rdx),
		R10: uintptr(regs.R10),
		R8:  uintptr(regs.R8),
		R9:  uintptr(regs.R9),
		Pid: this.cmd.Process.Pid,
	}); err != nil {
		return err
	}

	return nil
}

func (this *Process) GetSignalInfo() (*abi.SignalInfo, error) {
	return abi.GetSignalInfo(this.cmd.Process.Pid)
}

func (this *Process) WaitForStatus() (syscall.WaitStatus, error) {
	var status syscall.WaitStatus
	_, err := syscall.Wait4(this.cmd.Process.Pid, &status, 0, nil)
	return status, err
}

func (this *Process) getRegs() (syscall.PtraceRegs, error) {
	var regs syscall.PtraceRegs
	err := syscall.PtraceGetRegs(this.cmd.Process.Pid, &regs)
	return regs, err
}
