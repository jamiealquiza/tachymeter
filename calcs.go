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

	m.Times = m.Times[:m.TimesUsed]
	sort.Sort(m.Times)

	metrics := &Metrics{}
	metrics.Samples = m.TimesUsed
	metrics.Count = m.Count
	metrics.Time.Total = calcTotal(m.Times)
	metrics.Time.Avg = calcAvg(metrics.Time.Total, metrics.Samples)
	metrics.Time.p95 = calcp95(m.Times)
	metrics.Time.Long5p = calcLong5p(m.Times)
	metrics.Time.Short5p = calcShort5p(m.Times)
	metrics.Time.Max = m.Times[metrics.Samples-1]
	metrics.Time.Min = m.Times[0]
	rateTime := float64(metrics.Samples) / float64(metrics.Time.Total)
	metrics.Rate.Second = rateTime * 1e9

	return metrics
}

// These should be self-explanatory:

func calcTotal(d []time.Duration) time.Duration {
	var t time.Duration
	for _, d := range d {
		t += d
	}

	return t
}

func calcAvg(d time.Duration, c int) time.Duration {
	return time.Duration(int(d) / c)
}

func calcp95(d []time.Duration) time.Duration {
	return d[int(float64(len(d))*0.95)]
}

func calcLong5p(d []time.Duration) time.Duration {
	set := d[int(float64(len(d))*0.95):]
	if len(set) == 0 {
		return time.Millisecond*0
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
	set := d[:int(float64(len(d))*0.05)]
	if len(set) == 0 {
		return time.Millisecond*0
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}
