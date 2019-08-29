This example measures the durations of each iteration of a single loop and prints the summarized output.

### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example/tachymeter-simple`


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

	// Measure events.
	for i := 0; i < 100; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		c.AddTime(time.Since(start))
	}

	// Calc output.
	results := c.Calc()

	// Print JSON format to console.
	fmt.Printf("%s\n\n", results.JSON())

	// Print pre-formatted console output.
	fmt.Printf("%s\n\n", results)

	// Print text histogram.
	fmt.Println(results.Histogram.String(15))
}
```

### Output
```
$ $GOPATH/bin/tachymeter-simple
{"Time":{"Cumulative":"671.871ms","HMean":"125.38µs","Avg":"13.43742ms","P50":"13.165ms","P75":"20.058ms","P95":"27.536ms","P99":"30.043ms","P999":"30.043ms","Long5p":"29.749ms","Short5p":"399.666µs","Max":"30.043ms","Min":"4µs","Range":"30.039ms","StdDev":"8.385117ms"},"Rate":{"Second":74.41904770409796},"Samples":50,"Count":100,"Histogram":[{"4µs - 3.007ms":5},{"3.007ms - 6.011ms":4},{"6.011ms - 9.015ms":10},{"9.015ms - 12.019ms":6},{"12.019ms - 15.023ms":7},{"15.023ms - 18.027ms":3},{"18.027ms - 21.031ms":4},{"21.031ms - 24.035ms":3},{"24.035ms - 27.039ms":3},{"27.039ms - 30.043ms":5}]}

50 samples of 100 events
Cumulative:	671.871ms
HMean:		125.38µs
Avg.:		13.43742ms
p50: 		13.165ms
p75:		20.058ms
p95:		27.536ms
p99:		30.043ms
p999:		30.043ms
Long 5%:	29.749ms
Short 5%:	399.666µs
Max:		30.043ms
Min:		4µs
Range:		30.039ms
StdDev:		8.385117ms
Rate/sec.:	74.42

       4µs - 3.007ms -----
   3.007ms - 6.011ms ---
   6.011ms - 9.015ms ---------------
  9.015ms - 12.019ms -------
 12.019ms - 15.023ms ---------
 15.023ms - 18.027ms -
 18.027ms - 21.031ms ---
 21.031ms - 24.035ms -
 24.035ms - 27.039ms -
 27.039ms - 30.043ms -----
```
