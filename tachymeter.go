// Package tachymeter yields summarized data
// describing a series of timed events.
package tachymeter

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Config holds tachymeter initialization parameters.
type Config struct {
	// Size defines the sample capacity.
	Size int

	// HBins is the number of histogram bins.
	HBins int

	// Safe is retained so that existing users keep compiling.
	//
	// Deprecated: Tachymeter is always safe for concurrent use.
	Safe bool
}

// Tachymeter holds event durations
// and counts.
type Tachymeter struct {
	size     uint64
	times    []atomic.Int64 // Ring buffer of event durations.
	count    atomic.Uint64
	wallTime atomic.Int64
	hBins    int
	mu       sync.Mutex // Serializes Calc and Reset.
}

// timeSlice holds time.Duration values.
type timeSlice []time.Duration

// Histogram is an ordered list of bins, where each bin is a
// map["low - high duration"]count of events that fall within
// the low-high time duration range.
type Histogram []map[string]uint64

// Metrics holds the calculated outputs
// produced from a Tachymeter sample set.
type Metrics struct {
	Time struct { // All values under Time are selected entirely from events within the sample window.
		Cumulative time.Duration // Cumulative time of all sampled events.
		HMean      time.Duration // Event duration harmonic mean.
		Avg        time.Duration // Event duration average.
		P50        time.Duration // Event duration nth percentiles ..
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Long5p     time.Duration // Average of the longest 5% event durations.
		Short5p    time.Duration // Average of the shortest 5% event durations.
		Max        time.Duration // Highest event duration.
		Min        time.Duration // Lowest event duration.
		StdDev     time.Duration // Standard deviation.
		Range      time.Duration // Event duration range (Max-Min).
	}
	Rate struct {
		// Per-second rate based on event duration avg. via Metrics.Cumulative / Metrics.Samples.
		// If SetWallTime was called, event duration avg = wall time / Metrics.Count
		Second float64
	}
	Histogram        *Histogram    // Frequency distribution of event durations in len(Histogram) bins of HistogramBinSize.
	HistogramBinSize time.Duration // The width of a histogram bin in time.
	Samples          int           // Number of events included in the sample set.
	Count            int           // Total number of events observed.
}

// New initializes a new Tachymeter. A sample window
// size below 1 is raised to 1; a histogram bin count
// below 1 falls back to the default of 10.
func New(c *Config) *Tachymeter {
	size := c.Size
	if size < 1 {
		size = 1
	}

	hBins := c.HBins
	if hBins < 1 {
		hBins = 10
	}

	return &Tachymeter{
		size:  uint64(size),
		times: make([]atomic.Int64, size),
		hBins: hBins,
	}
}

// Reset resets a Tachymeter instance for reuse,
// clearing the event count and any wall time set
// with SetWallTime.
func (m *Tachymeter) Reset() {
	m.mu.Lock()
	m.count.Store(0)
	m.wallTime.Store(0)
	m.mu.Unlock()
}

// AddTime adds a time.Duration to Tachymeter.
func (m *Tachymeter) AddTime(t time.Duration) {
	m.times[(m.count.Add(1)-1)%m.size].Store(int64(t))
}

// SetWallTime optionally sets an elapsed wall time duration.
// This affects rate output by using total events counted over time.
// This is useful for concurrent/parallelized events that overlap
// in wall time and are writing to a shared Tachymeter instance.
func (m *Tachymeter) SetWallTime(t time.Duration) {
	m.wallTime.Store(int64(t))
}

// WriteHTML writes a histograph
// html file to the cwd.
func (m *Metrics) WriteHTML(p string) error {
	w := Timeline{}
	w.AddEvent(m)
	return w.WriteHTML(p)
}

// String satisfies the String interface.
func (m *Metrics) String() string {
	return fmt.Sprintf(`%d samples of %d events
Cumulative:	%s
HMean:		%s
Avg.:		%s
p50: 		%s
p75:		%s
p95:		%s
p99:		%s
p999:		%s
Long 5%%:	%s
Short 5%%:	%s
Max:		%s
Min:		%s
Range:		%s
StdDev:		%s
Rate/sec.:	%.2f`,
		m.Samples,
		m.Count,
		m.Time.Cumulative,
		m.Time.HMean,
		m.Time.Avg,
		m.Time.P50,
		m.Time.P75,
		m.Time.P95,
		m.Time.P99,
		m.Time.P999,
		m.Time.Long5p,
		m.Time.Short5p,
		m.Time.Max,
		m.Time.Min,
		m.Time.Range,
		m.Time.StdDev,
		m.Rate.Second)
}

// JSON returns a *Metrics as
// a JSON string.
func (m *Metrics) JSON() string {
	j, _ := json.Marshal(m)

	return string(j)
}

// MarshalJSON defines the output formatting
// for the JSON() method. This is exported as a
// requirement but not intended for end users.
func (m *Metrics) MarshalJSON() ([]byte, error) {
	// Durations are rendered as their
	// human-readable strings.
	type times struct {
		Cumulative string
		HMean      string
		Avg        string
		P50        string
		P75        string
		P95        string
		P99        string
		P999       string
		Long5p     string
		Short5p    string
		Max        string
		Min        string
		Range      string
		StdDev     string
	}

	return json.Marshal(&struct {
		Time      times
		Rate      struct{ Second float64 }
		Samples   int
		Count     int
		Histogram *Histogram
	}{
		Time: times{
			Cumulative: m.Time.Cumulative.String(),
			HMean:      m.Time.HMean.String(),
			Avg:        m.Time.Avg.String(),
			P50:        m.Time.P50.String(),
			P75:        m.Time.P75.String(),
			P95:        m.Time.P95.String(),
			P99:        m.Time.P99.String(),
			P999:       m.Time.P999.String(),
			Long5p:     m.Time.Long5p.String(),
			Short5p:    m.Time.Short5p.String(),
			Max:        m.Time.Max.String(),
			Min:        m.Time.Min.String(),
			Range:      m.Time.Range.String(),
			StdDev:     m.Time.StdDev.String(),
		},
		Rate:      struct{ Second float64 }{m.Rate.Second},
		Samples:   m.Samples,
		Count:     m.Count,
		Histogram: m.Histogram,
	})
}

// String returns a formatted Histogram string with
// bar lengths scaled to a width of s.
func (h *Histogram) String(s int) string {
	if h == nil || len(*h) == 0 {
		return ""
	}

	// Get the histogram min/max counts.
	var low, high uint64 = math.MaxUint64, 0
	for _, bin := range *h {
		for _, v := range bin {
			low = min(low, v)
			high = max(high, v)
		}
	}

	// With a single bin, scale its
	// bar to the full width.
	if len(*h) == 1 {
		low = 0
	}

	// Build the histogram string.
	var b strings.Builder
	for _, bin := range *h {
		for k, v := range bin {
			blen := scale(float64(v), float64(low), float64(high), 1, float64(s))
			fmt.Fprintf(&b, "%22s %s\n", k, strings.Repeat("-", int(blen)))
		}
	}

	return b.String()
}

// scale maps the input x in the range [a0, a1]
// to the output range [b0, b1].
func scale(x, a0, a1, b0, b1 float64) float64 {
	if x == a0 {
		return b0
	}

	return (x-a0)/(a1-a0)*(b1-b0) + b0
}
