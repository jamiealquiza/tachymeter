This example uses the tachymeter [Timeline](https://godoc.org/github.com/jamiealquiza/tachymeter#Timeline) type for gathering summary metrics from several iterations of a measured loop, outputting a single HTML page with histograms and summaries per iteration.

### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example/tachymeter-graphing`


### Example code
```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})
	tl := tachymeter.Timeline{}

	// Run 3 iterations of a loop that we're
	// interesting in summarizing with tachymeter.
	for iter := 0; iter <3; iter++ {
		fmt.Printf("Running iteration %d\n", iter)
		// Capture timing data from the loop.
		for i := 0; i < 100; i++ {
			start := time.Now()
			time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
			c.AddTime(time.Since(start))
		}

		// Add each loop tachymeter
		// to the event timeline.
		tl.AddEvent(c.Calc())
		c.Reset()
	}


	// Write out an HTML page with the
	// histogram for all iterations.
	err := tl.WriteHTML(".")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Results written")
}
```

### Output
```
$ $GOPATH/bin/tachymeter-graphing
Running iteration 0
Running iteration 1
Running iteration 2
Results written
```

## [HTML Example Output](https://jamiealquiza.github.io/tachymeter/)
