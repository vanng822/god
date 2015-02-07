# god

Keep an eye on some process running

# Usage
	
### build
	
	go build
	
### run

	>> ./god --pidfile example.pid -s ./example/test_bin  -p 8080 -c config.conf
	
### restart

	>> kill -s HUP $(cat example.pid)
	
### stop
	
	>> kill $(cat example.pid)