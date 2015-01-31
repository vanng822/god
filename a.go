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

func (a *Args) Parse() {
	args := os.Args[1:]
	max := len(args)

	for i := 0; i < max; i++ {
		if args[i] == "-pid" {
			i++
			if !isArgValue(args[i]) {
				panic("Invalid pid value")
			}
			a.pidFile = args[i]
			continue
		}
		if args[i] == "-s" {
			i++
			if !isArgValue(args[i]) {
				panic("Invalid program")
			}
			a.programs = append(a.programs, args[i])
			continue
		}
		if args[i] == "-f" {
			a.force = true
			continue
		}
		a.args = append(a.args, args[i])
	}
}

func isArgValue(value string) bool {
	return strings.Index(value, "-") != 0
}
