package god

import (
	"fmt"
	"log"
	"strings"
)

type Args struct {
	args     []string
	pidFile  string
	force    bool
	programs []string
	help     bool
}

func usage() {
	log.Printf(`
Usage: go run daemonize.go --pidfile daemonize.pid --pidclean -s program args ...
	--pidfile   A pidfile which process id will be stored.
	--pidclean  Clean old pidfile if there is one and try to run this program
	-s          Program you want to run in background. Can repeat for multiple daemons
	args        Other arguments that you want to pass to daemons
	
"go run daemonize.go" is your daemonize program. See example/main.go at https://github.com/vanng822/god

Example: go run example/main.go --pidfile god.pid -s ./example/test_bin -p 8080
	`)
}

func (a *Args) Parse(args []string) error {
	max := len(args)
	for i := 0; i < max; i++ {
		if args[i] == "--help" {
			a.help = true
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
		if args[i] == "-s" {
			i++
			if i >= max || !isArgValue(args[i]) {
				return fmt.Errorf("Invalid program")
			}
			a.programs = append(a.programs, args[i])
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

func isArgValue(value string) bool {
	return strings.Index(value, "-") != 0
}
