/*
 * The MIT License
 *
 *  Copyright (c) 2016, Yoann Dubreuil
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 *
 */
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	// Version is a state variable, written at the link stage. See Makefile.
	Version string
	// CommitID is a state variable, written at the link stage. See Makefile.
	CommitID string
)

func usage() {
	fmt.Println(`usage: trampoline <subcommand>

Where subcommand can be:
	cdexec: Run a command after changing current working directory to the given directory
	wait: wait to be killed, no more, no less!
`)
}

func cdexec(args []string) {
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	err := os.Chdir(args[0])
	if err != nil {
		log.Fatal(err)
		os.Exit(255)
	}

	binary, lookErr := exec.LookPath(args[1])
	if lookErr != nil {
		log.Fatal(lookErr)
		os.Exit(255)
	}

	if err := syscall.Exec(binary, args[1:], os.Environ()); err != nil {
		log.Fatal(err)
		os.Exit(255)
	}
}

func wait_system_signal() {
	channel := make(chan os.Signal)
	signal.Notify(channel)
	<-channel
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	subCommand := os.Args[1]
	commandLine := os.Args[2:]

	switch subCommand {
	case "cdexec":
		cdexec(commandLine)
	case "wait":
		wait_system_signal()
	default:
		usage()
	}
}
