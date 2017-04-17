[![GoDoc](https://godoc.org/github.com/jamiealquiza/tachymeter?status.svg)](https://godoc.org/github.com/jamiealquiza/tachymeter)

# tachymeter

Tachymeter simplifies the process of creating summarized rate and latency information from a series of timed events: _"In a loop with 1,000 database calls, what was the 95%ile and lowest observed latency? What was the per-second rate?"_

Event durations in the form of [`time.Duration`](https://golang.org/pkg/time/#Duration) are added to a tachymeter instance using the `AddTime(t time.Duration)` method.

After all desired timings have been gathered, the `Calc()` is called, returning a [`*Metrics`](https://godoc.org/github.com/jamiealquiza/tachymeter#Metrics). The results held by a `*Metrics` can be accessed in several ways (where `t` represents a tachymeter instance):

- `tachymeter.Metrics` for direct access:
```golang
results := t.Calc
fmt.Printf("Median latency: %s\n", results.Time.P50)
```

 - JSON string: `jsonResults := c.Json()`
 ```golang
results := t.Json()
fmt.Printf("%s\n\n", results)
```
 - Printing pre-formatted output to console:
 ```golang
 t.results.Dump()`
 ```
 
 - HTML histogram (written to the absolute path specified as the method parameter)
 ```golang
 err := t.Calc().WriteHtml(".")
 ```

# Usage

Tachymeter is initialized with a Size parameter that specifies the max sample size that will be used in the calculation. This is done to control resource usage and minimise the impact of introducting tachymeter into your application (by favoring fixed-size slices with indexed inserts rather than appends, limiting sort times, etc.); the `AddTime` method is o(1) @ ~20ns on modern hardware. If your actual event count is smaller than or equal to the tachymeter sample size, all of the meaused events will be included. If the event count exceeds the tachymeter size, the oldest data will be overwritten (this results in a last-window sample; sampling configuration will eventually be added).

See the [example](https://github.com/jamiealquiza/tachymeter/tree/master/example) file for a fully functioning example.

```golang
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
Cumulative:     705.24222ms
Avg.:           14.104844ms
p50:            13.073198ms
p75:            21.358238ms
p95:            28.289403ms
p99:            30.544326ms
p999:           30.544326ms
Long 5%:        29.843555ms
Short 5%:       356.145µs
Max:            30.544326ms
Min:            2.455µs
Range:          30.541871ms
Rate/sec.:      70.90
```

### Output Descriptions

- `Cumulative`: Aggregate of all sample durations.
- `Avg.`: Average event duration per sample.
- `p<N>`: Nth %ile.
- `Long 5%`: Average event duration of the longest 5%.
- `Short 5%`: Average event duration of the shortest 5%.
- `Max`: Max observed event duration.
- `Min`: Min observed event duration.
- `Range`: The delta between the max and min sample time
- `Rate/sec.`: Per-second rate based on cumulative time and sample count.

# Accurate Rates With Parallelism

By default, tachymeter calcualtes rate based on the number of events possible per-second according to average event duration. This model doesn't work in asynchronous or parallelized scenarios since events may be overlapping in time. For example, with many Goroutines writing durations to a shared tachymeter in parallel, the global rate must be determined by using the total event count over the total wall time elapsed.

Tachymeter exposes a `SetWallTime` method for these scenarios.

Example:

```golang
<...>

func main() {
    // Initialize tachymeter.
    c := tachymeter.New(&tachymeter.Config{Size: 50})

    // Start wall time for all Goroutines.
    wallTimeStart := time.Now()
    var wg sync.WaitGroup
    
    // Run tasks.
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go someTask(t, wg)
    }
    
    wg.Wait()

    // When finished, set elapsed wall time.
    t.SetWallTime(time.Since(wallTimeStart))
    
    // Rate outputs will be accurate.
    fmt.Println(t.Calc().Dump())
}

func someTask(t *tachymeter.Tachymeter, wg *sync.WaitGroup) {
    defer wg.Done()
    start := time.Now()

    // Task we're timing added here.
    t.AddTime(time.Since(start))
}

<...>
```
