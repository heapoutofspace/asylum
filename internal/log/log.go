package log

import (
	"fmt"
	"os"
)

const (
	blue   = "\033[34m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)

func Info(format string, args ...any) {
	fmt.Printf(blue+"i"+reset+" "+format+"\n", args...)
}

func Success(format string, args ...any) {
	fmt.Printf(green+"ok"+reset+" "+format+"\n", args...)
}

func Warn(format string, args ...any) {
	fmt.Printf(yellow+"!"+reset+" "+format+"\n", args...)
}

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, red+"x"+reset+" "+format+"\n", args...)
}

func Build(format string, args ...any) {
	fmt.Printf(cyan+"#"+reset+" "+format+"\n", args...)
}
