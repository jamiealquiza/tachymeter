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
	fmt.Printf("%s\n\n", results.String())

	// Print text histogram.
	fmt.Println(results.Histogram.String(15))
}
```

### Output
```
$ $GOPATH/bin/tachymeter-simple
{"Time":{"Cumulative":"694.643892ms","HMean":"59.157µs","Avg":"13.892877ms","P50":"11.457787ms","P75":"20.14303ms","P95":"29.169708ms","P99":"31.149213ms","P999":"31.149213ms","Long5p":"30.62519ms","Short5p":"429.969µs","Max":"31.149213ms","Min":"1.778µs","Range":"31.147435ms"},"Rate":{"Second":71.97932721475654},"Samples":50,"Count":100,"Histogram":[{"1µs - 3.116ms":5},{"3.116ms - 6.231ms":5},{"6.231ms - 9.346ms":8},{"9.346ms - 12.46ms":8},{"12.46ms - 15.575ms":5},{"15.575ms - 18.69ms":4},{"18.69ms - 21.804ms":4},{"21.804ms - 24.919ms":4},{"24.919ms - 28.034ms":2},{"28.034ms - 31.149ms":5}]}

50 samples of 100 events
Cumulative:     694.643892ms
HMean:          59.157µs
Avg.:           13.892877ms
p50:            11.457787ms
p75:            20.14303ms
p95:            29.169708ms
p99:            31.149213ms
p999:           31.149213ms
Long 5%:        30.62519ms
Short 5%:       429.969µs
Max:            31.149213ms
Min:            1.778µs
Range:          31.147435ms
Rate/sec.:      71.98

       1µs - 3.116ms --------
   3.116ms - 6.231ms --------
   6.231ms - 9.346ms ---------------
   9.346ms - 12.46ms ---------------
  12.46ms - 15.575ms --------
  15.575ms - 18.69ms -----
  18.69ms - 21.804ms -----
 21.804ms - 24.919ms -----
 24.919ms - 28.034ms -
 28.034ms - 31.149ms --------
```
