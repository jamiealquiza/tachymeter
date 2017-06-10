package tachymeter

import (
	"fmt"
	"math"
	"sort"
	"sync/atomic"
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

// Calc summarizes Tachymeter sample data
// and returns it in the form of a *Metrics.
func (m *Tachymeter) Calc() *Metrics {
	metrics := &Metrics{}
	if atomic.LoadUint64(&m.Count) == 0 {
		return metrics
	}

	m.Lock()

	metrics.Samples = int(math.Min(float64(atomic.LoadUint64(&m.Count)), float64(m.Size)))
	metrics.Count = int(atomic.LoadUint64(&m.Count))
	times := make(timeSlice, metrics.Samples)
	copy(times, m.Times[:metrics.Samples])
	sort.Sort(times)

	metrics.Time.Cumulative = calcTimeCumulative(times)
	var rateTime float64
	if m.WallTime != 0 {
		rateTime = float64(metrics.Count) / float64(m.WallTime)
	} else {
		rateTime = float64(metrics.Samples) / float64(metrics.Time.Cumulative)
	}

	metrics.Rate.Second = rateTime * 1e9

	m.Unlock()

	metrics.Time.Avg = calcAvg(times)
	metrics.Time.HMean = calcHMean(times)
	metrics.Time.P50 = times[len(times)/2]
	metrics.Time.P75 = calcP(times, 0.75)
	metrics.Time.P95 = calcP(times, 0.95)
	metrics.Time.P99 = calcP(times, 0.99)
	metrics.Time.P999 = calcP(times, 0.999)
	metrics.Time.Long5p = calcLong5p(times)
	metrics.Time.Short5p = calcShort5p(times)
	metrics.Time.Max = times[metrics.Samples-1]
	metrics.Time.Min = times[0]
	metrics.Time.Range = metrics.Time.Max - metrics.Time.Min

	var bSize time.Duration
	metrics.Histogram, bSize = calcHgram(m.HBuckets, times, metrics.Time.Min, metrics.Time.Max, metrics.Time.Range)
	metrics.HistogramBucketSize = bSize

	return metrics
}

// calcHgram returns a histogram of event durations t in b buckets,
// along with the bucket size.
// A histogram bucket is a map["low-high duration"]count of events that
// fall within the low / high range.
func calcHgram(b int, t timeSlice, low, max, r time.Duration) ([]map[string]int, time.Duration) {
	// Interval is the time range / n buckets.
	interval := time.Duration(int64(r) / int64(b))
	high := low + interval
	hgram := []map[string]int{}
	pos := 1 // Bucket position.

	bstring := fmt.Sprintf("%s - %s", low, high)
	bucket := map[string]int{}

	for _, v := range t {
		// If v fits in the current bucket,
		// increment the bucket count.
		if v <= high {
			bucket[bstring]++
		} else {
			// If not, prepare the next bucket.
			hgram = append(hgram, bucket)
			bucket = map[string]int{}

			// Update the high/low range values.
			low = high + time.Nanosecond

			high += interval
			// if we're going into the
			// last bucket, set high to max.
			if pos == b-1 {
				high = max
			}

			bstring = fmt.Sprintf("%s - %s", low, high)

			// The value didn't fit in the previous
			// bucket, so the new bucket count should
			// be incremented.
			bucket[bstring]++

			pos++
		}
	}

	hgram = append(hgram, bucket)

	return hgram, interval
}

// These should be self-explanatory:

func calcHMean(d []time.Duration) time.Duration {
	var total float64

	for _, t := range d {
		total += (1 / float64(t))
	}

	return time.Duration(float64(len(d)) / total)
}

func calcTimeCumulative(d []time.Duration) time.Duration {
	var total time.Duration
	for _, t := range d {
		total += t
	}

	return total
}

func calcAvg(d []time.Duration) time.Duration {
	var total time.Duration
	for _, t := range d {
		total += t
	}
	return time.Duration(int(total) / len(d))
}

func calcP(d []time.Duration, p float64) time.Duration {
	return d[int(float64(len(d))*p+0.5)-1]
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
