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

	for i := 0; i < 100; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		c.AddTime(time.Since(start))
	}

	resultsj := c.Json()
	fmt.Printf("%s\n\n", resultsj)

	// Print the pre-formatted console output.
	results.Calc().Dump()
}
```

### Output
```
$ $GOPATH/bin/tachymeter-simple
{"Time":{"Cumulative":"669.072395ms","Avg":"13.381447ms","P50":"12.257547ms","P75":"19.807712ms","P95":"27.389804ms","P99":"31.873821ms","P999":"31.873821ms","Long5p":"30.012355ms","Short5p":"439.781µs","Max":"31.873821ms","Min":"2.018µs","Range":"31.871803ms"},"Rate":{"Second":74.73032869634383},"Samples":50,"Count":100,"Histogram":[{"2.018µs-3.189198ms":6},{"3.189199ms - 6.376378ms":2},{"6.376379ms - 9.563558ms":12},{"9.563559ms - 12.750738ms":6},{"12.750739ms - 15.937918ms":7},{"15.937919ms - 19.125098ms":4},{"19.125099ms - 22.312278ms":2},{"22.312279ms - 25.499458ms":6},{"25.499459ms - 28.686638ms":4},{"28.686639ms - 31.873821ms":1}]}

50 samples of 100 events
Cumulative:     669.072395ms
Avg.:           13.381447ms
p50:            12.257547ms
p75:            19.807712ms
p95:            27.389804ms
p99:            31.873821ms
p999:           31.873821ms
Long 5%:        30.012355ms
Short 5%:       439.781µs
Max:            31.873821ms
Min:            2.018µs
Range:          31.871803ms
Rate/sec.:      74.73
```
