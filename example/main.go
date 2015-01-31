package main

import (
	"github.com/vanng822/god"
	"runtime"
	"log"
)

func main() {
	log.Printf("Main:Number of goroutines %d", runtime.NumGoroutine())
	z := god.NewGoz()
	z.Start()
	log.Printf("Main:Number of goroutines %d", runtime.NumGoroutine())
}