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
	sort.Sort(m.Times)

	metrics := &Metrics{}
	metrics.Samples = m.TimesUsed
	metrics.Count = m.Count
	metrics.Time.Total = calcTotal(m.Times)
	metrics.Time.Avg = calcAvg(metrics.Time.Total, len(m.Times))
	metrics.Time.p95 = calcp95(m.Times)
	metrics.Time.Long10p = calcLong10p(m.Times)
	metrics.Time.Short10p = calcShort10p(m.Times)
	metrics.Time.Max = m.Times[len(m.Times)-1]
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
	return d[int(float64(len(d))*0.9)]
}

func calcLong10p(d []time.Duration) time.Duration {
	var t time.Duration
	var i int
	for _, n := range d[int(float64(len(d))*0.9):] {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func calcShort10p(d []time.Duration) time.Duration {
	var t time.Duration
	var i int
	for _, n := range d[:int(float64(len(d))*0.1)] {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}
