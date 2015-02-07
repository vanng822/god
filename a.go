package main

import (
	"fmt"
	"strconv"
	"strings"
)

func usage() {
	fmt.Printf(`
Usage: god --pidfile daemonize.pid --pidclean --interval 86400 -s program args ...
	--pidfile   A pidfile which process id will be stored.
	--pidclean  Clean old pidfile if there is one and try to run this program
	--interval  Number of seconds to restart all programs, given in integer. Minimum is %.0f seconds
	-s          Program you want to run in background. Can repeat for multiple daemons
	args        Other arguments that you want to pass to daemons
	
	god --version For printing out version

Example: god --pidfile god.pid -s sleep 10` + "\n", MIMIMUM_AGE)
}

type Args struct {
	args        []string
	pidFile     string
	force       bool
	programs    []string
	programArgs [][]string
	help        bool
	version     bool
	interval    int
}

func (a *Args) Parse(args []string) error {
	max := len(args)
	for i := 0; i < max; i++ {
		if args[i] == "--help" {
			a.help = true
			return nil
		}
		if args[i] == "--version" {
			a.version = true
			return nil
		}
		if args[i] == "--pidfile" {
			i++
			if i >= max || !isArgValue(args[i]) {
				return fmt.Errorf("Invalid pidfile value")
			}
			a.pidFile = args[i]
			continue
		}
		if args[i] == "--interval" {
			i++
			if i >= max || !isArgValue(args[i]) {
				return fmt.Errorf("Invalid interval value")
			}
			interval, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("Invalid interval value")
			}
			if float64(interval) <= MIMIMUM_AGE {
				return fmt.Errorf("Minium value for interval is %.0f seconds", MIMIMUM_AGE)
			}
			a.interval = interval
			continue
		}
		if args[i] == "-s" {
			i++
			if i >= max || !isArgValue(args[i]) {
				return fmt.Errorf("Invalid program")
			}
			a.programs = append(a.programs, args[i])
			programArgs := findProgramArgs(args, i+1)
			a.programArgs = append(a.programArgs, programArgs)
			i += len(programArgs)
			continue
		}
		if args[i] == "--pidclean" {
			a.force = true
			continue
		}
		a.args = append(a.args, args[i])
	}

	return nil
}

func findProgramArgs(args []string, start int) []string {
	var pargs []string
	max := len(args)
	for n := start; n < max; n++ {
		if !isArgValue(args[n]) {
			break
		}
		pargs = append(pargs, args[n])
	}

	return pargs
}

func isArgValue(value string) bool {
	return strings.Index(value, "-") != 0
}
