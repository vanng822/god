package main

import (
	"flag"
	"os"
	"fmt"
)
func main() {
	help := flag.Bool("help", false, "Show usage")
	version := flag.Bool("version", false, "Show version and exit")
	
	flag.Parse()
	if *help {
		usage()
		os.Exit(0)
	}
	if *version {
		fmt.Printf("god %s\n", VERSION)
		os.Exit(0)
	}
	z := NewGoz()
	z.Start()
}