package tachymeter

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkAddTime(b *testing.B) {
	b.StopTimer()

	ta := New(&Config{Size: b.N})
	d := time.Millisecond

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.AddTime(d)
	}
}

func BenchmarkAddTimeSampling(b *testing.B) {
	b.StopTimer()

	ta := New(&Config{Size: 100})
	d := time.Millisecond

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.AddTime(d)
	}
}

func TestNew(t *testing.T) {
	// Out-of-range configuration values
	// fall back to sane defaults.
	ta := New(&Config{Size: -1, HBins: -1})

	if len(ta.times) != 1 {
		t.Errorf("Expected sample window size 1, got %d", len(ta.times))
	}

	if ta.hBins != 10 {
		t.Errorf("Expected 10 histogram bins, got %d", ta.hBins)
	}
}

func TestReset(t *testing.T) {
	ta := New(&Config{Size: 3})

	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.SetWallTime(time.Second)
	ta.Reset()

	if ta.count.Load() != 0 {
		t.Errorf("Expected count 0, got %d", ta.count.Load())
	}

	if ta.wallTime.Load() != 0 {
		t.Errorf("Expected wall time 0, got %d", ta.wallTime.Load())
	}
}

func TestAddTime(t *testing.T) {
	ta := New(&Config{Size: 3})

	ta.AddTime(time.Millisecond)

	if time.Duration(ta.times[0].Load()) != time.Millisecond {
		t.Fail()
	}
}

func TestAddTimeWraps(t *testing.T) {
	ta := New(&Config{Size: 3})

	// The fourth event overwrites
	// the oldest slot.
	for i := 1; i <= 4; i++ {
		ta.AddTime(time.Duration(i) * time.Millisecond)
	}

	expected := []time.Duration{
		4 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
	}

	for i, d := range expected {
		if got := time.Duration(ta.times[i].Load()); got != d {
			t.Errorf("Expected %s at slot %d, got %s", d, i, got)
		}
	}
}

func TestSetWallTime(t *testing.T) {
	ta := New(&Config{Size: 3})

	ta.SetWallTime(time.Millisecond)

	if time.Duration(ta.wallTime.Load()) != time.Millisecond {
		t.Fail()
	}
}

// TestConcurrent exercises concurrent AddTime, Calc, and
// Reset calls; it exists to be run with the race detector.
func TestConcurrent(t *testing.T) {
	ta := New(&Config{Size: 64})

	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				ta.AddTime(time.Duration(i) * time.Microsecond)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = ta.Calc()
		}
		ta.Reset()
	}()

	wg.Wait()
}
