package main

import (
	"rundc/pkg/kernel"
	"rundc/rundc"

	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		kernel := kernel.Kernel{}
		kernel.Init()
		if len(os.Args) < 3 { //for debug
			kernel.Run("echo", []string{"hello"})
		}
		kernel.Run(os.Args[2], os.Args[3:])
	default:
		cli := &rundc.Cli{}
		cli.Main(os.Args)
	}

}
