package log

import (
	"fmt"
	"os"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func Info(message interface{}) {
	fmt.Println(green, fmt.Sprintf("[INFO] %s", message), reset)
}

func Error(message interface{}) {
	fmt.Println(red, fmt.Sprintf("[ERROR] %s", message), reset)
}

func ErrorAndExit(message interface{}) {
	Error(message)
	os.Exit(1)
}
