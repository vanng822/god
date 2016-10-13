# god

Keep an eye on some process running. God will restart your program if it exits unexpected. This is helpful when you don't have time 24/7 to watch over your applications. It restarts program on SIGHUP and forward SIGUSR1 and SIGUSR2.

Be aware this may not work well if your program forks another process, special in watch mode.

# Usage
	
### build
	
	go build
	
### run

	>> ./god --pidfile god.pid -s go run test_program/test_bin.go

### Check test_bin.go working
	
	//Open in browser
	http://127.0.0.1:8080/
	
	
### run in watch mode, for a go program, ie don't run "go run"

	>> ./god --watch folder1,folder2 --watch-exts go,json --pidfile god.pid -s make build-go-program
	>> ./god --watch touchfolder --watch-exts touch --pidfile god2.pid -s /path/to/go-program


### restart

	>> kill -s HUP $(cat god.pid)
	
### stop
	
	>> kill $(cat god.pid)
