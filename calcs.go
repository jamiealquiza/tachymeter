package tachymeter

import (
	"sort"
	"time"
)

// Satisfy sort for timeSlice.
// Sorts in increasing order of duration.

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return int64(p[i]) < int64(p[j])
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Calc calcs data held in a *Tachymeter
// and returns a *Metrics.
func (m *Tachymeter) Calc() *Metrics {
	m.Lock()
	defer m.Unlock()

	metrics := &Metrics{}
	if m.Count == 0 {
		return metrics
	}

	times := m.Times[:m.TimesUsed]
	sort.Sort(times)

	metrics.Samples = m.TimesUsed
	metrics.Count = m.Count
	metrics.Time.Total = calcTimeTotal(times)
	metrics.Time.Avg = calcAvg(times, metrics.Samples)
	metrics.Time.Median = times[len(times)/2]
	metrics.Time.P95 = calcP95(times)
	metrics.Time.Long5p = calcLong5p(times)
	metrics.Time.Short5p = calcShort5p(times)
	metrics.Time.Max = times[metrics.Samples-1]
	metrics.Time.Min = times[0]

	var rateTime float64
	if m.WallTime != 0 {
		rateTime = float64(metrics.Count) / float64(m.WallTime)
	} else {
		rateTime = float64(metrics.Samples) / float64(metrics.Time.Total)
	}
	metrics.Rate.Second = rateTime * 1e9

	return metrics
}

// These should be self-explanatory:

func calcTimeTotal(d []time.Duration) time.Duration {
	var total time.Duration
	for _, t := range d {
		total += t
	}

	return total
}

func calcAvg(d []time.Duration, c int) time.Duration {
	var total time.Duration
	for _, t := range d {
		total += t
	}
	return time.Duration(int(total) / c)
}

func calcP95(d []time.Duration) time.Duration {
	return d[int(float64(len(d))*0.95+0.5)-1]
}

func calcLong5p(d []time.Duration) time.Duration {
	set := d[int(float64(len(d))*0.95+0.5):]

	if len(set) <= 1 {
		return d[len(d)-1]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func calcShort5p(d []time.Duration) time.Duration {
	set := d[:int(float64(len(d))*0.05+0.5)]

	if len(set) <= 1 {
		return d[0]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}
