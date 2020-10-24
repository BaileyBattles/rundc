package process

import (
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
	cmd *exec.Cmd
}

func CreatePtraceProcess(path string) *Process {
	cmd := exec.Command(path)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	return &Process{
		cmd: cmd,
	}
}
