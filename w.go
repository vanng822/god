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
	stat, err := folder.Stat()
	if err != nil {
		return nil, err
	}
	// path can be a file
	if !stat.IsDir() {
		return folders, nil
	}
	defer folder.Close()
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
		if err = watcher.Add(path); err != nil {
			log.Fatal(err)
		}
		folders, err := scan(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range folders {
			if err = watcher.Add(f); err != nil {
				log.Fatal(err)
			}
		}
	}
	allExts := strings.Split(exts, ",")

	for {
		select {
		case event := <-watcher.Events:
			if len(allExts) > 0 {
				for _, ext := range allExts {
					if strings.HasSuffix(event.Name, ext) {
						restart <- true
						break
					}
				}
			} else {
				restart <- true
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
