package kernel

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"rundc/pkg/abi"
	"rundc/pkg/sys"
)

type Process struct {
	kernel *Kernel
	cmd    *exec.Cmd
}

func CreatePtraceProcess(path string, args []string) *Process {
	fmt.Println("Creating process")
	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	kernel := &Kernel{
		Table: NewSyscallTable(),
	}
	kernel.Table.Register(uintptr(1), func(SyscallArguments) (uintptr, uintptr, syscall.Errno) {
		fmt.Println("you do what i say")
		return uintptr(0), uintptr(0), syscall.Errno(uintptr(1))
	})
	return &Process{
		cmd:    cmd,
		kernel: kernel,
	}

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
	sys.PrintSyscallName(regs.Orig_rax)

	this.kernel.HandleSyscall(this, uintptr(regs.Orig_rax), SyscallArguments{
		Rdi: uintptr(regs.Rdi),
		Rsi: uintptr(regs.Rsi),
		Rdx: uintptr(regs.Rdx),
		R10: uintptr(regs.R10),
		R8:  uintptr(regs.R8),
		R9:  uintptr(regs.R9),
	})
	// regs, err = this.getRegs()
	// if err != nil {
	// 	return err
	// }
	//abi.SetReturnValue(this.cmd.Process.Pid, int(rval))
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
