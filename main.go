package main

import (
	"os/exec"
	"rundc/pkg/kernel"
	"rundc/pkg/log"
	"rundc/rundc"
	"syscall"

	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		kernel := kernel.Kernel{}
		kernel.Init()
		kernel.Run(os.Args[2], os.Args[3:])
	case "run2":
		cmd := exec.Command("ls", "")
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS,
		}
		err := cmd.Start()
		if err != nil {
			log.ErrorAndExit(err)
		}
	default:
		cli := &rundc.Cli{}
		cli.Main(os.Args)
	}

}
