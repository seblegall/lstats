# lstats
A golang load tests package

## Requirements
* `github.com/gosuri/uiprogress` package

## Installation

Get the package using `go get` : 

```sh
go get -v github.com/seblegall/lstats
```

## Usage

Example of usage :

```golang
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/seblegall/lstats"
)

func main() {

	content, _ := ioutil.ReadFile("urls.txt")
	urls := strings.Split(string(content), "\n")

	var reqs []*http.Request

	for _, url := range urls {
		if url != "" {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("Error when creatin request for url : %s", url)
			}
			req.SetBasicAuth("test", "test")

			reqs = append(reqs, req)
		}
	}

	//25 is the total parallel workers desired.
	//This way, lstat will do the calls on 25 differents go routine.
	test := lstats.NewLoadStats(reqs, 25)
	test.Launch()
	test.Print()
}
```

*note* : The `NewLoadStats()` function is expecting a slice of `*http.Request`. This let you create custom request and add headers, auth, or anything you need to actualy to the request.
The second parameters is the count of parrallel workers desired.

Now, let's create an `urls.txt` file containing a list of url to test with one url by line : 

```
http://example.com/test/123/tests
http://example.com/test/456/tests
http://example.com/test/789/tests
http://example.com/test/101/tests
http://example.com/test/121/tests
http://example.com/test/314/tests
http://example.com/test/151/tests
```

and finaly, we can launch the test using : 

```sh
$ go run main.go
```

This will output something similar to : 
```
Starting load test...
  23s [===================================================================>] 100%

AVG Response Time       Total Calls in error
0.578685                0
```


## TODO

- [ ] Add unit tests
- [ ] Add possibiity to set the number of concurrent workers
- [ ] Add a random wait time between requests in order to simulate more realistic user calls