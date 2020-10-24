package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"rundc/pkg/abi"
	"rundc/pkg/sys"
)

type Process struct {
	cmd          *exec.Cmd
	syscallStart bool
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
	return &Process{
		cmd:          cmd,
		syscallStart: true,
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
		fmt.Printf("Ptrace failed with errno = %s\n", errno.Error())
		return error(errno)
	}
	return nil
}

func (this *Process) HandleSyscall(doIt bool) error {
	if doIt {
		abi.PutRegs(this.cmd.Process.Pid)
	}
	if !this.syscallStart {
		this.syscallStart = true
		return nil
	}
	this.syscallStart = false
	regs, err := this.getRegs()
	if err != nil {
		return err
	}
	sys.PrintSyscallName(regs.Orig_rax)
	return nil
}

func (this *Process) GetSignalInfo() (*abi.SignalInfo, error) {
	// sigInfo := &sys.SignalInfo{}
	// ptr := unsafe.Pointer(sigInfo)
	// if _, _, errno := syscall.Syscall6(syscall.SYS_PTRACE, syscall.PTRACE_GETSIGINFO, uintptr(this.cmd.Process.Pid), 0, uintptr(ptr), 0, 0); errno != 0 {
	// 	fmt.Printf("Ptrace failed with errno = %s\n", errno.Error())
	// 	return nil, error(errno)
	// }
	// return sigInfo, nil
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
