# god

Keep an eye on some process running

# Usage
	
### main.go

	import (
		"github.com/vanng822/god"
	)
	func main() {
		z := god.NewGoz()
		z.Start()
	}
	
### run

	>> go run main.go -pid example.pid -s ./example/test_bin  -p 8080 -c config.conf
	
### restart

	>> kill -s HUP $(cat example.pid)
	
### stop
	
	>> kill $(cat example.pid)