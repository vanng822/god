# god

Keep an eye on some process running

# Usage
	
### build
	
	go build
	
### run

	>> ./god --pidfile god.pid -s go run test_program/test_bin.go

### Check test_bin.go working
	
	//Open in browser
	http://127.0.0.1:8080/

### restart

	>> kill -s HUP $(cat example.pid)
	
### stop
	
	>> kill $(cat example.pid)