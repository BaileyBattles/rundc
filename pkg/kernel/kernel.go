package kernel

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"rundc/pkg/abi"
	"rundc/pkg/log"
	"rundc/pkg/kernel/syscalls"
	"rundc/pkg/sys"
)

type Kernel struct {
	table *SyscallTable
}

func (kernel *Kernel) Init() {
	kernel.table = NewSyscallTable()
}

func (kernel *Kernel) HandleSyscall(process *Process, id uintptr, args syscalls.SyscallArguments) error {
	if syscall := kernel.table.GetSyscall(id); syscall != nil {
		if _, _, errno := syscall(args); errno != 0 {
			fmt.Printf("Failed executing syscall with errno: %s\n", errno.Error())
			return errno
		}
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
	return nil
}

func (this *Kernel) createPtraceProcess(path string, args []string) *Process {
	log.Info("Creating process")
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
	go this.runProcess(p, ch)
	<-ch
	kernelLoop()
}

func (this *Kernel) runProcess(p *Process, ch chan struct{}) {
	for {
		p.Ptrace(sys.PTRACE_SYSCALL)

		status, err := p.WaitForStatus()
		if status.Exited() {
			log.Info("Child process has exited")
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
	select {}
}
