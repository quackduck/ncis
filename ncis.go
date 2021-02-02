package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	version = "dev"
	helpMsg = `lolsh - A shell with all output lolcat-ed
Usage: lolsh [-v/--version | -v/--help]`
)

func main() {
	if len(os.Args) > 4 {
		handleErrStr("too many arguments (>3)") // will accept ip addr, pass list and user list
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("help", "h"); hasOption {
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("version", "v"); hasOption {
		fmt.Println("ncis " + version)
		return
	}
}

func argsHaveOption(long string, short string) (hasOption bool, foundAt int) {
	for i, arg := range os.Args {
		if arg == "--"+long || arg == "-"+short {
			return true, i
		}
	}
	return false, 0
}

func handleErr(err error) {
	handleErrStr(err.Error())
}

func handleErrStr(str string) {
	_, _ = fmt.Fprintln(os.Stderr, color.RedString("error: ")+str)
}
