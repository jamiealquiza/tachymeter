package tachymeter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Size int
	Safe bool // Optionally lock if concurrent access is needed.
}

type timeSlice []time.Duration

type Tachymeter struct {
	sync.Mutex
	Safe          bool
	TimeTotal     time.Duration
	Times         timeSlice
	TimesPosition int
	TimesUsed     int
	Count         int
	WallTime      time.Duration
}

type Metrics struct {
	Time struct {
		Total   time.Duration
		Avg     time.Duration
		Median  time.Duration
		p95     time.Duration
		Long5p  time.Duration
		Short5p time.Duration
		Max     time.Duration
		Min     time.Duration
	}
	Rate struct {
		Second float64
	}
	Samples int
	Count   int
}

// New initializes a new tachymeter.
func New(c *Config) *Tachymeter {
	return &Tachymeter{
		Times: make([]time.Duration, c.Size),
		Safe:  c.Safe,
	}
}

// Reset resets the counters / positions.
func (m *Tachymeter) Reset() {
	m.Lock()
	defer m.Unlock()

	m.TimesPosition, m.TimesUsed, m.Count = 0, 0, 0
}

// AddTime adds a time.Duration to tachymeter.
func (m *Tachymeter) AddTime(t time.Duration) {
	if m.Safe {
		m.Lock()
		defer m.Unlock()
	}

	// Add to total time.
	m.TimeTotal += t

	// If we're at the end, rollover and
	// start overwriting.
	if m.TimesPosition == len(m.Times) {
		m.TimesPosition = 0
	}

	m.Times[m.TimesPosition] = t
	m.TimesPosition++
	if m.TimesUsed < len(m.Times) {
		m.TimesUsed++
	}
}

// AddCount increments the tachymeter event
// count by i.
func (m *Tachymeter) AddCount(i int) {
	if m.Safe {
		m.Lock()
		defer m.Unlock()
	}

	m.Count += i
}

func (m *Tachymeter) SetWallTime(t time.Duration) {
	if m.Safe {
		m.Lock()
		defer m.Unlock()
	}

	m.WallTime = t
}

// Dump prints a formatted summary of tachymeter metrics.
func (m *Metrics) Dump() {
	fmt.Printf("%d samples of %d events\n", m.Samples, m.Count)
	fmt.Printf("Total:\t\t%s\n", m.Time.Total)
	fmt.Printf("Avg.:\t\t%s\n", m.Time.Avg)
	fmt.Printf("Median: \t%s\n", m.Time.Median)
	fmt.Printf("95%%ile:\t\t%s\n", m.Time.p95)
	fmt.Printf("Longest 5%%:\t%s\n", m.Time.Long5p)
	fmt.Printf("Shortest 5%%:\t%s\n", m.Time.Short5p)
	fmt.Printf("Max:\t\t%s\n", m.Time.Max)
	fmt.Printf("Min:\t\t%s\n", m.Time.Min)
	fmt.Printf("Rate/sec.:\t%.2f\n", m.Rate.Second)
}

// Json returns a json string of tachymeter metrics.
func (m *Tachymeter) Json() string {
	metrics := m.Calc()
	j, _ := json.Marshal(&metrics)

	return string(j)
}

// MarshalJSON defines the output formatting
// for the Json() method.
func (m *Metrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Time struct {
			Total   string
			Avg     string
			Median  string
			p95     string
			Long5p  string
			Short5p string
			Max     string
			Min     string
		}
		Rate struct {
			Second float64
		}
		Samples int
		Count   int
	}{
		Time: struct {
			Total   string
			Avg     string
			Median  string
			p95     string
			Long5p  string
			Short5p string
			Max     string
			Min     string
		}{
			Total:   m.Time.Total.String(),
			Avg:     m.Time.Avg.String(),
			Median:  m.Time.Median.String(),
			p95:     m.Time.p95.String(),
			Long5p:  m.Time.Long5p.String(),
			Short5p: m.Time.Short5p.String(),
			Max:     m.Time.Max.String(),
			Min:     m.Time.Min.String(),
		},
		Rate: struct{ Second float64 }{
			Second: m.Rate.Second,
		},
		Samples: m.Samples,
		Count:   m.Count,
	})
}
