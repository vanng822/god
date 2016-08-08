package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchDirs(dirs, exts string, restart chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if dirs == "" {
		dirs = "."
	}

	allDirs := strings.Split(dirs, ",")
	log.Println("Watching dirs:", allDirs)
	for _, dd := range allDirs {
		path, err := filepath.Abs(dd)
		if err != nil {
			log.Fatal(err)
		}
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}
	allExts := strings.Split(exts, ",")
	var shouldRestart bool
	for {
		select {
		case event := <-watcher.Events:
			//log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				for _, ext := range allExts {
					if strings.HasSuffix(event.Name, ext) {
						shouldRestart = true
						log.Println("modified file:", event.Name, exts)
						restart <- shouldRestart
						break
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
