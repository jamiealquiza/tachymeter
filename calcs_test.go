package tachymeter_test

import (
	//"fmt"
	"sort"
	"testing"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func TestCalc(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 30})

	// These end up overwritten; we're
	// putting 32 events in a size 30 Tachymeter.
	ta.AddTime(12 * time.Millisecond)
	ta.AddTime(96 * time.Millisecond)

	ta.AddTime(9 * time.Millisecond)
	ta.AddTime(4 * time.Millisecond)
	ta.AddTime(88 * time.Millisecond)
	ta.AddTime(37 * time.Millisecond)
	ta.AddTime(42 * time.Millisecond)
	ta.AddTime(77 * time.Millisecond)
	ta.AddTime(93 * time.Millisecond)
	ta.AddTime(89 * time.Millisecond)
	ta.AddTime(12 * time.Millisecond)
	ta.AddTime(36 * time.Millisecond)
	ta.AddTime(54 * time.Millisecond)
	ta.AddTime(21 * time.Millisecond)
	ta.AddTime(17 * time.Millisecond)
	ta.AddTime(14 * time.Millisecond)
	ta.AddTime(67 * time.Millisecond)

	ta.AddTime(9 * time.Millisecond)
	ta.AddTime(4 * time.Millisecond)
	ta.AddTime(88 * time.Millisecond)
	ta.AddTime(37 * time.Millisecond)
	ta.AddTime(42 * time.Millisecond)
	ta.AddTime(77 * time.Millisecond)
	ta.AddTime(93 * time.Millisecond)
	ta.AddTime(89 * time.Millisecond)
	ta.AddTime(12 * time.Millisecond)
	ta.AddTime(36 * time.Millisecond)
	ta.AddTime(54 * time.Millisecond)
	ta.AddTime(21 * time.Millisecond)
	ta.AddTime(17 * time.Millisecond)
	ta.AddTime(14 * time.Millisecond)
	ta.AddTime(67 * time.Millisecond)

	metrics := ta.Calc()

	if metrics.Samples != 30 {
		t.Error("Expected 30, got ", metrics.Samples)
	}

	if metrics.Count != 32 {
		t.Error("Expected 32, got ", metrics.Count)
	}

	expectedDurs := []time.Duration{
		4000000,
		4000000,
		9000000,
		9000000,
		12000000,
		12000000,
		14000000,
		14000000,
		17000000,
		17000000,
		21000000,
		21000000,
		36000000,
		36000000,
		37000000,
		37000000,
		42000000,
		42000000,
		54000000,
		54000000,
		67000000,
		67000000,
		77000000,
		77000000,
		88000000,
		88000000,
		89000000,
		89000000,
		93000000,
		93000000,
	}

	sort.Sort(ta.Times)

	for n, d := range ta.Times {
		if d != expectedDurs[n] {
			t.Errorf("Expected %d, got %d\n", expectedDurs[n], d)
		}
	}

	if metrics.Time.Cumulative != 1320000000 {
		t.Errorf("Expected 1320000000, got %d\n", metrics.Time.Cumulative)
	}

	if metrics.Time.Avg != 44000000 {
		t.Errorf("Expected 44000000, got %d\n", metrics.Time.Avg)
	}

	if metrics.Time.P50 != 37000000 {
		t.Errorf("Expected 37000000, got %d\n", metrics.Time.P50)
	}

	if metrics.Time.P75 != 77000000 {
		t.Errorf("Expected 77000000, got %d\n", metrics.Time.P75)
	}

	if metrics.Time.P95 != 93000000 {
		t.Errorf("Expected 93000000, got %d\n", metrics.Time.P95)
	}

	if metrics.Time.P99 != 93000000 {
		t.Errorf("Expected 93000000, got %d\n", metrics.Time.P99)
	}

	if metrics.Time.P999 != 93000000 {
		t.Errorf("Expected 93000000, got %d\n", metrics.Time.P999)
	}

	if metrics.Time.Long5p != 93000000 {
		t.Errorf("Expected 93000000, got %d\n", metrics.Time.Long5p)
	}

	if metrics.Time.Short5p != 4000000 {
		t.Errorf("Expected 4000000, got %d\n", metrics.Time.Short5p)
	}

	if metrics.Time.Max != 93000000 {
		t.Errorf("Expected 93000000, got %d\n", metrics.Time.Max)
	}

	if metrics.Time.Min != 4000000 {
		t.Errorf("Expected 4000000, got %d\n", metrics.Time.Min)
	}

	if metrics.Time.Range != 89000000 {
		t.Errorf("Expected 89000000, got %d\n", metrics.Time.Range)
	}

	if metrics.Rate.Second != 22.72727272727273 {
		t.Errorf("Expected 22.73, got %0.2f\n", metrics.Rate.Second)
	}
}
