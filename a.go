package god

import (
	"os"
	"strings"
)

type Args struct {
	args     []string
	pidFile  string
	force    bool
	programs []string
}

func NewArgs() *Args {
	args := new(Args)
	args.args = make([]string, 0)
	return args
}

func (as *Args) isValue(value string) bool {
	return strings.Index(value, "-") != 0
}

func (as *Args) Parse() error {
	args := os.Args[1:]
	max := len(args)

	for i := 0; i < max; i++ {
		if args[i] == "-pid" {
			i++
			if !as.isValue(args[i]) {
				panic("Invalid pid value")
			}
			as.pidFile = args[i]
			continue
		}
		if args[i] == "-s" {
			i++
			if !as.isValue(args[i]) {
				panic("Invalid program")
			}
			as.programs = append(as.programs, args[i])
			continue
		}
		if args[i] == "-f" {
			as.force = true
			continue
		}
		as.args = append(as.args, args[i])
	}
	return nil
}
