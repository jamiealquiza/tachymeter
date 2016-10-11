# tachymeter

Tachymeter simplifies the process of gathering summarized latency and rate statistics from a series of work. 

After initializing a new `tachymeter`, latencies in the form of [`time.Duration`](https://golang.org/pkg/time/#Duration) that measure an event duration are added with the `AddTime()` function. Event counts, usually one per event duration, are added with the `AddCount()` function. Event counts may not correlate 1:1 with the number of event durations depending on what information is desired. For example, timing a function call and incrementing by 1 will yield latency and rates regarding making those function calls, whereas timing a single function call but incrementing by the *number of results returned* can be used to infer a different piece of information.

Tachymeter is initialized with a Size parameter that specifies the max sample size that will be used in the calculation. This is done to control resource usage and minimise the impact of introducting tachymeter into your application (by avoiding slice appends, reducing sort times, etc.). If your actual event count is smaller than the tachymeter sample size, 100% of your data will be included. If the actual event count exceeds the tachymeter size, the oldest data will be overwritten.

After all of the desired latencies have been gathered, tachymeter data can be gathered or viewed in several ways:
 - Raw data accessible as a `tachymeter.Metrics`: `results := c.Calc`
 - A json string: `jsonResults := c.Json()`
 - Printing a pre-formatted output to console: `resultsc.Dump()`

# Example Usage

See the [example](https://github.com/jamiealquiza/tachymeter/tree/master/example) file for a fully functioning example.

```go
import "github.com/jamiealquiza/tachymeter"

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})

	for i := 0; i < 100; i++ {
		start := time.Now()
		doSomeWork()
		c.AddTime(time.Since(start))
		c.AddCount(1)
	}

	c.Calc().Dump()
}
```

```
50 samples of 100 events
Total:			746.317344ms
Avg.:			14.926346ms
Median: 		13.902036ms
95%ile:			30.33816ms
Longest 5%:		31.541162ms
Shortest 5%:	466.619µs
Max:			32.146993ms
Min:			3.444µs
Rate/sec.:		67.00
```
