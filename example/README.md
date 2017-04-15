### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example`


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

	for i := 0; i < 100; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		c.AddTime(time.Since(start))
	}

	resultsj := c.Json()
	fmt.Printf("%s\n\n", resultsj)

	results := c.Calc()
	// Print the pre-formatted console output.
	results.Dump()
	// Create an HTML graph of the event histogram.
	results.DumpHistogramGraph()
}
```

### Output
```
$ $GOPATH/bin/example
{"Time":{"Cumulative":"700.215421ms","Avg":"14.004308ms","P50":"12.990493ms","P75":"20.061121ms","P95":"30.612424ms","P99":"31.925942ms","P999":"31.925942ms","Long5p":"31.87406ms","Short5p":"465.319µs","Max":"31.925942ms","Min":"2.8µs","Range":"31.923142ms"},"Rate":{"Second":71.40659645654962},"Samples":50,"Count":100,"Histogram":[{"2.8µs-3.195114ms":5},{"3.195115ms - 6.387428ms":4},{"6.387429ms - 9.579742ms":9},{"9.579743ms - 12.772056ms":7},{"12.772057ms - 15.96437ms":7},{"15.964371ms - 19.156684ms":3},{"19.156685ms - 22.348998ms":5},{"22.348999ms - 25.541312ms":3},{"25.541313ms - 28.733626ms":3},{"28.733627ms - 31.92594ms":3},{"31.925941ms - 31.925942ms":1}]}

50 samples of 100 events
Cumulative:     700.215421ms
Avg.:           14.004308ms
p50:            12.990493ms
p75:            20.061121ms
p95:            30.612424ms
p99:            31.925942ms
p999:           31.925942ms
Long 5%:        31.87406ms
Short 5%:       465.319µs
Max:            31.925942ms
Min:            2.8µs
Range:          31.923142ms
Rate/sec.:      71.41
```

# HTML histogram output
![ss](https://cloud.githubusercontent.com/assets/4108044/25065346/a7f23398-21cb-11e7-8d39-c97a1ea4d136.png)
