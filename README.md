[![GoDoc](https://godoc.org/github.com/jamiealquiza/tachymeter?status.svg)](https://godoc.org/github.com/jamiealquiza/tachymeter)

# tachymeter

Tachymeter simplifies the process of gathering summarized rate and latency information from a series of timed events: _"In a loop with 1,000 database calls, what was the 95%ile and lowest observed latency? What was the per-second rate?"_

# Usage

After initializing a `tachymeter`, event durations in the form of [`time.Duration`](https://golang.org/pkg/time/#Duration) are added using the `AddTime(t time.Duration)` method. Once all desired timing have been collected, the data is summarized by calling the `Calc()` method (returning a [`*Metrics`](https://godoc.org/github.com/jamiealquiza/tachymeter#Metrics)). `*Metrics` fields can be accessed directly or via several [output methods](https://github.com/jamiealquiza/tachymeter#output-options).

See further code [examples](https://github.com/jamiealquiza/tachymeter/tree/master/example).

```golang
import "github.com/jamiealquiza/tachymeter"

func main() {
    // Initialize a tachymeter with a max
    // sample window of 50.
    t := tachymeter.New(&tachymeter.Config{Size: 50})

    for i := 0; i < 100; i++ {
        start := time.Now()
        doSomeWork()
        // We add the time that
        // each doSomeWork() call took.
        t.AddTime(time.Since(start))
    }

    // The timing summaries are calculated
    // and dumped to console.
    t.Calc().Dump()
}
```

```
50 samples of 100 events
Cumulative:     669.433006ms
HMean:          47.477µs
Avg.:           13.38866ms
p50:            11.191119ms
p75:            19.15929ms
p95:            28.145686ms
p99:            30.135862ms
p999:           30.135862ms
Long 5%:        29.156558ms
Short 5%:       424.823µs
Max:            30.135862ms
Min:            1.765µs
Range:          30.134097ms
Rate/sec.:      74.69
```

### Output Descriptions

- `Cumulative`: Aggregate of all sample durations.
- `HMean`: Event duration harmonic mean.
- `Avg.`: Average event duration per sample.
- `p<N>`: Nth %ile.
- `Long 5%`: Average event duration of the longest 5%.
- `Short 5%`: Average event duration of the shortest 5%.
- `Max`: Max observed event duration.
- `Min`: Min observed event duration.
- `Range`: The delta between the max and min sample time
- `Rate/sec.`: Per-second rate based on cumulative time and sample count.

### Configuration

Tachymeter is initialized with a `Size` parameter that specifies the max sample count that can be held. This is done to control resource usage and minimize the impact of tachymeter inside an application; the `AddTime` method is o(1) @ ~20ns on modern hardware. If the actual event count is smaller than or equal to the configured tachymeter size, all of the meaused events will be included in the calculated results. If the event count exceeds the tachymeter size, the oldest data will be overwritten. In this scenario, the last window of data (that fits into the configured `Size`) will be used for output calculations.

# Output Options

`t` represents a tachymeter instance.

### `tachymeter.Metrics` for direct access
```golang
results := t.Calc()
fmt.Printf("Median latency: %s\n", results.Time.P50)
```

### JSON string
 ```golang
fmt.Printf("%s\n\", results.JSON())
```
### Printing pre-formatted output to console
 ```golang
results.Dump()
 ```

### HTML histogram
 Tachymeter `*Metrics` results also have to ability to be written as HTML histograms. The `WriteHTML(p string)` method is called where `p` is an output path where the HTML file should be written.

 ```golang
 err := results.WriteHTML(".")
 ```
 
![ss](https://cloud.githubusercontent.com/assets/4108044/25826873/c40d62b8-3405-11e7-9dec-047d1e0c6f42.png)

Tachymeter also provides a `Timeline` type that's used to gather a series of `*Metrics` (each `*Metrics` themselves holding data summarizing a series of measured events). `*Metrics` are added to a `*Timeline` using the `AddEvent(m *Metrics)` method. Once the desired number of `*Metrics` has been collected, `WriteHTML` can be called on the `*Timeline`, resulting in an single HTML page with a histogram for each captured `*Metrics`. An example use case may be a benchmark where tachymeter is used to summarize the timing results of a loop, but several iterations of the loop are to be called in series. See the [tachymeter-graphing example](https://github.com/jamiealquiza/tachymeter/tree/master/example/tachymeter-graphing) for further details.

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
    
    // doSomeSlowDbCall()

    // Task we're timing added here.
    t.AddTime(time.Since(start))
}

<...>
```
