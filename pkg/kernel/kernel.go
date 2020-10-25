package kernel

import (
	"fmt"
	"os"
	"os/exec"
	"rundc/pkg/abi"
	"rundc/pkg/kernel/syscalls"
	"rundc/pkg/sys"
	"syscall"
)

type Kernel struct {
	table *SyscallTable
}

func (kernel *Kernel) Init() {
	kernel.table = NewSyscallTable()
}

func (kernel *Kernel) HandleSyscall(process *Process, id uintptr, args syscalls.SyscallArguments) {
	if syscall := kernel.table.GetSyscall(id); syscall != nil {
		syscall(args)
		abi.SetSyscall(process.cmd.Process.Pid, 64)
		err := waitForSyscallCompletion(process)
		if err != nil {
			panic("Failed to wait for get PiD")
		}
		abi.SetReturnValue(process.cmd.Process.Pid, int(args.Rdx))
	} else {
		err := waitForSyscallCompletion(process)
		if err != nil {
			panic("Failed to wait for syscall")
		}
	}
}

func (this *Kernel) createPtraceProcess(path string, args []string) *Process {
	fmt.Println("Creating process")
	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	return &Process{
		cmd:    cmd,
		kernel: this,
	}

}

func (this *Kernel) Run(path string, args []string) {
	p := this.createPtraceProcess(path, args)
	p.Start()
	_, err := p.WaitForStatus()
	if err != nil {
		fmt.Printf("Wait returned with err: %v\n\n\n", err.Error())
	}
	ch := make(chan struct{})
	this.runProcess(p, ch)
	<-ch
	kernelLoop()
}

func (this *Kernel) runProcess(p *Process, ch chan struct{}) {
	for {
		p.Ptrace(sys.PTRACE_SYSCALL)

		status, err := p.WaitForStatus()
		if status.Exited() {
			fmt.Println("Child process has exited")
			close(ch)
			return
		}
		if err != nil {
			fmt.Printf("Error waiting for status %s\n", err.Error())
		}

		if status.Stopped() && status.StopSignal() == syscall.SIGTRAP {

			err = p.HandleSyscall()
			if err != nil {
				fmt.Printf("Error handling syscall: %s\n", err.Error())
			}

		}
		if status.Stopped() && status.StopSignal() == syscall.SIGSEGV {
			panic("Received a signal")
		}
	}
}

func waitForSyscallCompletion(process *Process) error {
	err := process.Ptrace(sys.PTRACE_SYSCALL)
	process.WaitForStatus()
	if err != nil {
		panic("Failed")
	}
	return nil
}

func kernelLoop() {
	for {
	}
}
