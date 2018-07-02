package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/martinlindhe/notify"
)

func ShowVersion() {
	fmt.Println("v1.0.0")
}

func main() {
	var help, version bool
	var seconds int
	var kill bool
	var alert bool

	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&version, "v", false, "show version")

	flag.IntVar(&seconds, "s", 60, "detect slow command, if not finished after N seconds")
	flag.BoolVar(&kill, "k", false, "kill slow command")
	flag.BoolVar(&alert, "a", false, "notify and play alert sound")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: " + os.Args[0] + " [OPTIONS] COMMAND [ARG...]")
		fmt.Println()
		fmt.Println("Run command, maybe it is slow")
		fmt.Println()
		fmt.Println("Options:")
		flag.CommandLine.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()

	args := flag.Args()

	if help {
		flag.Usage()
		return
	}

	if version {
		ShowVersion()
		return
	}

	if len(args) <= 0 {
		flag.Usage()
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slow := false
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	runc := make(chan error)
	start := time.Now()
	go func() {
		runc <- cmd.Run()
	}()

	sigc := make(chan os.Signal)
	signal.Notify(sigc)

	for {
		select {
		case sig := <-sigc:
			cmd.Process.Signal(sig)

		case <-ticker.C:
			ticker.Stop()

			slow = true

			if kill {
				cmd.Process.Kill()
			}

		case err := <-runc:
			if slow {
				title := "finish: " + time.Now().Sub(start).String()
				if err != nil {
					title = "error: " + err.Error()
				}

				if alert {
					notify.Alert(os.Args[0], title, strings.Join(args, " "), "")
				} else {
					notify.Notify(os.Args[0], title, strings.Join(args, " "), "")
				}
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}

			return
		}
	}
}
