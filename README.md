# tachymeter

Tachymeter simplifies the process of measuring and summarizing rate and latency data from a series of events.

Latencies in the form of [`time.Duration`](https://golang.org/pkg/time/#Duration) that measure an event duration are added to a tachymeter instance using the `AddTime()` method.

After all desired latencies have been gathered, tachymeter data can be retrieved in several ways:
 - Raw data accessible as a `tachymeter.Metrics`: `results := c.Calc`
 - A json string: `jsonResults := c.Json()`
 - Printing a pre-formatted output to console: `results.Dump()`

Tachymeter is initialized with a Size parameter that specifies the max sample size that will be used in the calculation. This is done to control resource usage and minimise the impact of introducting tachymeter into your application (by favoring fixed-size slices with indexed inserts rather than appends, limiting sort times, etc.). If your actual event count is smaller than the tachymeter sample size, 100% of your data will be included. If the actual event count exceeds the tachymeter size, the oldest data will be overwritten (this results in a last-window sample; sampling configuration will eventually be added).



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
	}

	c.Calc().Dump()
}
```

```
50 samples of 100 events
Cumulative:	649.553805ms
Avg.:		12.991076ms
p50: 		12.062832ms
p75:		18.256032ms
p95:		27.175521ms
p99:		28.182451ms
p999:		28.182451ms
Long 5%:	28.170356ms
Short 5%:	416.962µs
Max:		28.182451ms
Min:		2.24µs
Rate/sec.:	76.98
```

### Output Descriptions

- `Cumulative`: Aggregate of all sample durations.
- `Avg.`: Average event duration per sample.
- `p<N>`: Nth %ile.
- `Long 5%`: Average event duration of the longest 5%.
- `Short 5%`: Average event duration of the shortest 5%.
- `Max`: Max observed event duration.
- `Min`: Min observed event duration.
- `Rate/sec.`: Per-second rate based on cumulative time and sample count.
