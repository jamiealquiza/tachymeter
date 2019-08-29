// Package tachymeter yields summarized data
// describing a series of timed events.
package tachymeter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Config holds tachymeter initialization
// parameters. Size defines the sample capacity.
// Tachymeter is thread safe.
type Config struct {
	Size  int
	Safe  bool // Deprecated. Flag held on to as to not break existing users.
	HBins int  // Histogram bins.
}

// Tachymeter holds event durations
// and counts.
type Tachymeter struct {
	sync.Mutex
	Size     uint64
	Times    timeSlice
	Count    uint64
	WallTime time.Duration
	HBins    int
}

// timeslice holds time.Duration values.
type timeSlice []time.Duration

// Satisfy sort for timeSlice.
func (p timeSlice) Len() int           { return len(p) }
func (p timeSlice) Less(i, j int) bool { return int64(p[i]) < int64(p[j]) }
func (p timeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Histogram is a map["low-high duration"]count of events that
// fall within the low-high time duration range.
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

// New initializes a new Tachymeter.
func New(c *Config) *Tachymeter {
	var hSize int
	if c.HBins != 0 {
		hSize = c.HBins
	} else {
		hSize = 10
	}

	return &Tachymeter{
		Size:  uint64(c.Size),
		Times: make([]time.Duration, c.Size),
		HBins: hSize,
	}
}

// Reset resets a Tachymeter
// instance for reuse.
func (m *Tachymeter) Reset() {
	// This lock is obviously not needed for
	// the m.Count update, rather to prevent a
	// Tachymeter reset while Calc is being called.
	m.Lock()
	atomic.StoreUint64(&m.Count, 0)
	m.Unlock()
}

// AddTime adds a time.Duration to Tachymeter.
func (m *Tachymeter) AddTime(t time.Duration) {
	m.Times[(atomic.AddUint64(&m.Count, 1)-1)%m.Size] = t
}

// SetWallTime optionally sets an elapsed wall time duration.
// This affects rate output by using total events counted over time.
// This is useful for concurrent/parallelized events that overlap
// in wall time and are writing to a shared Tachymeter instance.
func (m *Tachymeter) SetWallTime(t time.Duration) {
	m.WallTime = t
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
	return json.Marshal(&struct {
		Time struct {
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
		Rate struct {
			Second float64
		}
		Samples   int
		Count     int
		Histogram *Histogram
	}{
		Time: struct {
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
		}{
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
		Rate: struct{ Second float64 }{
			Second: m.Rate.Second,
		},
		Histogram: m.Histogram,
		Samples:   m.Samples,
		Count:     m.Count,
	})
}

// String returns a formatted Metrics string scaled
// to a width of s.
func (h *Histogram) String(s int) string {
	if h == nil {
		return ""
	}

	var min, max uint64 = math.MaxUint64, 0
	// Get the histogram min/max counts.
	for _, bin := range *h {
		for _, v := range bin {
			if v > max {
				max = v
			}
			if v < min {
				min = v
			}
		}
	}

	// Handle cases of no or
	// a single bin.
	switch len(*h) {
	case 0:
		return ""
	case 1:
		min = 0
	}

	var b bytes.Buffer

	// Build histogram string.
	for _, bin := range *h {
		for k, v := range bin {
			// Get the bar length.
			blen := scale(float64(v), float64(min), float64(max), 1, float64(s))
			line := fmt.Sprintf("%22s %s\n", k, strings.Repeat("-", int(blen)))
			b.WriteString(line)
		}
	}

	return b.String()
}

// Scale scales the input x with the input-min a0,
// input-max a1, output-min b0, and output-max b1.
func scale(x, a0, a1, b0, b1 float64) float64 {
	a, b := x-a0, a1-a0
	var c float64
	if a == 0 {
		c = 0
	} else {
		c = a / b
	}
	return c*(b1-b0) + b0
}
