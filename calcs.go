package tachymeter

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Calc summarizes Tachymeter sample data
// and returns it in the form of a *Metrics.
func (m *Tachymeter) Calc() *Metrics {
	metrics := &Metrics{}

	// Take a snapshot of the sample
	// window and counts.
	m.Lock()

	if m.Count == 0 {
		m.Unlock()
		return metrics
	}

	metrics.Count = int(m.Count)
	metrics.Samples = metrics.Count
	if m.Count > m.Size {
		metrics.Samples = int(m.Size)
	}

	times := make(timeSlice, metrics.Samples)
	copy(times, m.Times[:metrics.Samples])
	wallTime := m.WallTime

	m.Unlock()

	sort.Sort(times)

	metrics.Time.Cumulative = times.cumulative()

	var rateTime float64
	switch {
	case wallTime != 0:
		rateTime = float64(metrics.Count) / float64(wallTime)
	case metrics.Time.Cumulative != 0:
		rateTime = float64(metrics.Samples) / float64(metrics.Time.Cumulative)
	}

	metrics.Rate.Second = rateTime * 1e9

	metrics.Time.Avg = times.avg()
	metrics.Time.HMean = times.hMean()
	metrics.Time.P50 = times[times.Len()/2]
	metrics.Time.P75 = times.p(0.75)
	metrics.Time.P95 = times.p(0.95)
	metrics.Time.P99 = times.p(0.99)
	metrics.Time.P999 = times.p(0.999)
	metrics.Time.Long5p = times.long5p()
	metrics.Time.Short5p = times.short5p()
	metrics.Time.Min = times.min()
	metrics.Time.Max = times.max()
	metrics.Time.Range = times.srange()
	metrics.Time.StdDev = times.stdDev()

	metrics.Histogram, metrics.HistogramBinSize = times.hgram(m.HBins)

	return metrics
}

// hgram returns a histogram of event durations in
// b bins, along with the bin size.
func (ts timeSlice) hgram(b int) (*Histogram, time.Duration) {
	min, max := ts.min(), ts.max()

	// Interval is the time range / n bins. A zero
	// interval (all samples within b nanoseconds of
	// each other) is raised to 1ns so that bins
	// cover a non-zero range.
	interval := (max - min) / time.Duration(b)
	if interval == 0 {
		interval = time.Nanosecond
	}

	// Tally each event in the bin
	// covering its duration.
	counts := make([]uint64, b)
	for _, v := range ts {
		bin := int((v - min) / interval)
		// The max value lands on the top boundary
		// of the last bin; clamp it in.
		if bin > b-1 {
			bin = b - 1
		}
		counts[bin]++
	}

	// Label each bin with the duration range it
	// covers, truncated to microsecond resolution.
	res := time.Duration(1000)
	hgram := make(Histogram, b)
	for i := range hgram {
		low := min + time.Duration(i)*interval
		high := low + interval
		if i == b-1 {
			high = max
		}

		bstring := fmt.Sprintf("%s - %s", low/res*res, high/res*res)
		hgram[i] = map[string]uint64{bstring: counts[i]}
	}

	return &hgram, interval
}

// These should be self-explanatory:

func (ts timeSlice) cumulative() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}

	return total
}

func (ts timeSlice) hMean() time.Duration {
	var total float64

	for _, t := range ts {
		total += (1 / float64(t))
	}

	return time.Duration(float64(ts.Len()) / total)
}

func (ts timeSlice) avg() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}
	return time.Duration(int(total) / ts.Len())
}

func (ts timeSlice) p(p float64) time.Duration {
	return ts[int(float64(ts.Len())*p+0.5)-1]
}

func (ts timeSlice) stdDev() time.Duration {
	m := ts.avg()
	s := 0.00

	for _, t := range ts {
		s += math.Pow(float64(m-t), 2)
	}

	msq := s / float64(ts.Len())

	return time.Duration(math.Sqrt(msq))
}

func (ts timeSlice) long5p() time.Duration {
	set := ts[int(float64(ts.Len())*0.95+0.5):]

	if len(set) <= 1 {
		return ts[ts.Len()-1]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func (ts timeSlice) short5p() time.Duration {
	set := ts[:int(float64(ts.Len())*0.05+0.5)]

	if len(set) <= 1 {
		return ts[0]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func (ts timeSlice) min() time.Duration {
	return ts[0]
}

func (ts timeSlice) max() time.Duration {
	return ts[ts.Len()-1]
}

func (ts timeSlice) srange() time.Duration {
	return ts.max() - ts.min()
}
