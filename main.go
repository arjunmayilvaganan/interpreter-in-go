package main

import (
	"fmt"
	"github.com/arjunmayilvaganan/nibbl/repl"
	"os"
	"os/user"
	"runtime"
	"time"
)

const PROGRAM = "nibbl"
const VERSION = "0.1"

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Now().Format(time.RFC1123Z))
	fmt.Printf("\nHello, %s\n", currentUser.Name)
	fmt.Printf("Welcome to %s v%s on %s\n", PROGRAM, VERSION, runtime.GOOS)

	repl.Start(os.Stdin, os.Stdout)
}
