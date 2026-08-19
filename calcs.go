package tachymeter

import (
	"fmt"
	"math"
	"slices"
	"time"
)

// Calc summarizes Tachymeter sample data
// and returns it in the form of a *Metrics.
func (m *Tachymeter) Calc() *Metrics {
	metrics := &Metrics{}

	// Take a snapshot of the sample
	// window and counts.
	m.mu.Lock()

	count := m.count.Load()
	if count == 0 {
		m.mu.Unlock()
		return metrics
	}

	metrics.Count = int(count)
	metrics.Samples = int(min(count, m.size))

	times := make(timeSlice, metrics.Samples)
	for i := range times {
		times[i] = time.Duration(m.times[i].Load())
	}
	wallTime := time.Duration(m.wallTime.Load())

	m.mu.Unlock()

	slices.Sort(times)

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
	metrics.Time.P50 = times.p(0.50)
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

	metrics.Histogram, metrics.HistogramBinSize = times.hgram(m.hBins)

	return metrics
}

// hgram returns a histogram of event durations in
// b bins, along with the bin size.
func (ts timeSlice) hgram(b int) (*Histogram, time.Duration) {
	low, high := ts.min(), ts.max()

	// Interval is the time range / n bins. A zero
	// interval (all samples within b nanoseconds of
	// each other) is raised to 1ns so that bins
	// cover a non-zero range.
	interval := (high - low) / time.Duration(b)
	if interval == 0 {
		interval = time.Nanosecond
	}

	// Tally each event in the bin
	// covering its duration.
	counts := make([]uint64, b)
	for _, v := range ts {
		// The max value lands on the top boundary
		// of the last bin; clamp it in.
		bin := min(int((v-low)/interval), b-1)
		counts[bin]++
	}

	// Label each bin with the duration range it
	// covers, truncated to microsecond resolution.
	const res = time.Microsecond
	hgram := make(Histogram, b)
	for i := range hgram {
		binLow := low + time.Duration(i)*interval
		binHigh := binLow + interval
		if i == b-1 {
			binHigh = high
		}

		label := fmt.Sprintf("%s - %s", binLow/res*res, binHigh/res*res)
		hgram[i] = map[string]uint64{label: counts[i]}
	}

	return &hgram, interval
}

// cumulative returns the sum of all event durations.
func (ts timeSlice) cumulative() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}

	return total
}

// avg returns the arithmetic mean event duration.
func (ts timeSlice) avg() time.Duration {
	return ts.cumulative() / time.Duration(len(ts))
}

// hMean returns the harmonic mean event duration.
func (ts timeSlice) hMean() time.Duration {
	var total float64
	for _, t := range ts {
		total += 1 / float64(t)
	}

	return time.Duration(float64(len(ts)) / total)
}

// p returns the nearest-rank pth percentile
// of the sorted timeSlice.
func (ts timeSlice) p(p float64) time.Duration {
	return ts[int(float64(len(ts))*p+0.5)-1]
}

// stdDev returns the population standard
// deviation of event durations.
func (ts timeSlice) stdDev() time.Duration {
	mean := ts.avg()

	var sqSum float64
	for _, t := range ts {
		d := float64(t - mean)
		sqSum += d * d
	}

	return time.Duration(math.Sqrt(sqSum / float64(len(ts))))
}

// long5p returns the average of the
// longest 5% event durations.
func (ts timeSlice) long5p() time.Duration {
	set := ts[int(float64(len(ts))*0.95+0.5):]
	if len(set) <= 1 {
		return ts.max()
	}

	return set.avg()
}

// short5p returns the average of the
// shortest 5% event durations.
func (ts timeSlice) short5p() time.Duration {
	set := ts[:int(float64(len(ts))*0.05+0.5)]
	if len(set) <= 1 {
		return ts.min()
	}

	return set.avg()
}

func (ts timeSlice) min() time.Duration {
	return ts[0]
}

func (ts timeSlice) max() time.Duration {
	return ts[len(ts)-1]
}

// srange returns the range (max-min)
// of event durations.
func (ts timeSlice) srange() time.Duration {
	return ts.max() - ts.min()
}
