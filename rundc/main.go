package main

import (
	"rundc/pkg/kernel"

	"os"
)

func main() {
	kernel := kernel.Kernel{}
	kernel.Init()
	if len(os.Args) < 2 { //for debug
		kernel.Run("echo", []string{"hello"})
	}
	kernel.Run(os.Args[1], os.Args[2:])
}
