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

func LogInfo(message string) {
	fmt.Println(string(green), fmt.Sprintf("[INFO] %s", message), string(reset))
}

func LogError(message string) {
	fmt.Println(string(red), fmt.Sprintf("[ERROR] %s", message), string(reset))
}

func LogErrorAndExit(message string) {
	LogError(message)
	os.Exit(1)
}
