package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("test_bin starting")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
		log.Println("Hello")
	})
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	log.Printf("test_bin listening at %d", 8080)
	go http.ListenAndServe(":8080", nil)
	sig := <-sigc
	log.Printf("Got signal: %s", sig)
}
