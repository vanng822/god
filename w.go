package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func scan(path string) ([]string, error) {
	var folders []string
	folder, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := folder.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, fi := range files {
		// skip all dot files/folders
		if fi.Name()[0] == '.' {
			continue
		}
		if fi.IsDir() {
			folders = append(folders, path+"/"+fi.Name())
			subfolder, err := scan(path + "/" + fi.Name())
			if err != nil {
				panic(err)
			}
			folders = append(folders, subfolder...)
		}
	}
	return folders, nil
}

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
	for _, dd := range allDirs {
		path, err := filepath.Abs(dd)
		if err != nil {
			log.Fatal(err)
		}
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
		folders, err := scan(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range folders {
			err = watcher.Add(f)
			if err != nil {
				log.Fatal(err)
			}
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