package common

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

const (
	LINES  = "-----> "
	BLANKS = "       "
	INFO   = " for more info"
)

var (
	Red, Yellow, Green func(a ...interface{}) string
)

func init() {
	if runtime.GOOS == "windows" {
		Red = fmt.Sprint
		Yellow = fmt.Sprint
		Green = fmt.Sprint
	} else {
		Red = color.New(color.FgRed).SprintFunc()
		Yellow = color.New(color.FgYellow).SprintFunc()
		Green = color.New(color.FgGreen).SprintFunc()
	}
}
