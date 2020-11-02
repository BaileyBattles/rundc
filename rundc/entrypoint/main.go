package main

import (
	"fmt"
	"os"
	"rundc/rundc"
)

func main() {
	fmt.Println("here")
	cli := rundc.Cli{}
	cli.Main(os.Args)
}
