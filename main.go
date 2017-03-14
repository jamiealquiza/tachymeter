// The MIT License (MIT)
//
// Copyright (c) 2016 Jamie Alquiza
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package tachymeter

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Config holds tachymeter initialization
// parameters. Size defines the sample capacity.
// Tachymeter is thread safe.
type Config struct {
	Size int
	Safe bool // Deprecated. Flag held on to as to not break existing users.
}

// timeslice is used to hold time.Duration values.
type timeSlice []time.Duration

// Tachymeter provides methods to collect
// sample durations and produce summarized
// latecy / rate output.
type Tachymeter struct {
	sync.Mutex
	Size     uint64
	Times    timeSlice
	Count    uint64
	WallTime time.Duration
}

// Metrics holds the calculated outputs
// produced from a Tachymeter sample set.
type Metrics struct {
	Time struct {
		Cumulative time.Duration
		Avg        time.Duration
		P50        time.Duration
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Long5p     time.Duration
		Short5p    time.Duration
		Max        time.Duration
		Min        time.Duration
		Range      time.Duration
	}
	Rate struct {
		Second float64
	}
	Samples int
	Count   int
}

// New initializes a new Tachymeter.
func New(c *Config) *Tachymeter {
	return &Tachymeter{
		Size:  uint64(c.Size),
		Times: make([]time.Duration, c.Size),
	}
}

// Reset resets a Tachymeter
// instance for reuse.
func (m *Tachymeter) Reset() {
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

// Dump prints a formatted Metrics output to console.
func (m *Metrics) Dump() {
	fmt.Printf("%d samples of %d events\n", m.Samples, m.Count)
	fmt.Printf("Cumulative:\t%s\n", m.Time.Cumulative)
	fmt.Printf("Avg.:\t\t%s\n", m.Time.Avg)
	fmt.Printf("p50: \t\t%s\n", m.Time.P50)
	fmt.Printf("p75:\t\t%s\n", m.Time.P75)
	fmt.Printf("p95:\t\t%s\n", m.Time.P95)
	fmt.Printf("p99:\t\t%s\n", m.Time.P99)
	fmt.Printf("p999:\t\t%s\n", m.Time.P999)
	fmt.Printf("Long 5%%:\t%s\n", m.Time.Long5p)
	fmt.Printf("Short 5%%:\t%s\n", m.Time.Short5p)
	fmt.Printf("Max:\t\t%s\n", m.Time.Max)
	fmt.Printf("Min:\t\t%s\n", m.Time.Min)
	fmt.Printf("Range:\t\t%s\n", m.Time.Range)
	fmt.Printf("Rate/sec.:\t%.2f\n", m.Rate.Second)
}

// Json calls the Calc method on a Tachymeter
// instance and returns a json string of the output.
func (m *Tachymeter) Json() string {
	metrics := m.Calc()
	j, _ := json.Marshal(&metrics)

	return string(j)
}

// MarshalJSON defines the output formatting
// for the Json() method. This is exported as a
// requirement but not intended for end users.
func (m *Metrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Time struct {
			Cumulative string
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
		}
		Rate struct {
			Second float64
		}
		Samples int
		Count   int
	}{
		Time: struct {
			Cumulative string
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
		}{
			Cumulative: m.Time.Cumulative.String(),
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
		},
		Rate: struct{ Second float64 }{
			Second: m.Rate.Second,
		},
		Samples: m.Samples,
		Count:   m.Count,
	})
}
