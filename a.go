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
	--graceful  Restart program in sequence
	--interval  Number of seconds to restart all programs, given in integer. Minimum is %.0f seconds
	--watch 		dirs to watch for file changes
 	--watch-exts file extenstions to watch, default by watching all extenstions
	--delay number of seconds to countdown before restarting the program when any file is changed
	--log	Log filename, output from stdout of child process
	--log-err Log filename for errors, output from stderr of child process
	-s program [args]   Program you want to run in background. Can repeat for multiple daemons
	args        Other arguments that you want to pass to daemons

	god --version For printing out version

Example: god --pidfile god.pid -s sleep 10`+"\n", MIMIMUM_AGE)
}

type Args struct {
	args            []string
	pidFile         string
	graceful        bool
	logFile         string
	logFileErr      string
	force           bool
	programs        []string
	programArgs     [][]string
	help            bool
	version         bool
	interval        int
	fileWatched     bool
	fileWatchedDirs string
	fileWatchedExts string
	delaySecs       int
	fileWatchedDone chan bool
}

func (a *Args) Parse(args []string) error {
	max := len(args)
	a.delaySecs = 3
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
		if args[i] == "--graceful" {
			i++
			a.graceful = true
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
		if args[i] == "--watch" {
			a.fileWatched = true
			i++
			a.fileWatchedDirs = args[i]
			continue
		}
		if args[i] == "--watch-exts" {
			a.fileWatched = true
			i++
			a.fileWatchedExts = args[i]
			continue
		}

		if args[i] == "--delay" {
			i++
			delay, err := strconv.Atoi(args[i])
			if err != nil {
				return err
			}
			a.delaySecs = delay
			continue
		}

		if args[i] == "--log" {
			i++
			a.logFile = args[i]
			continue
		}

		if args[i] == "--log-err" {
			i++
			a.logFileErr = args[i]
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
		if isNextProgram(args[n]) {
			break
		}
		pargs = append(pargs, args[n])
	}

	return pargs
}

func isNextProgram(arg string) bool {
	return arg == "-s"
}

func isArgValue(value string) bool {
	return strings.Index(value, "-") != 0
}
