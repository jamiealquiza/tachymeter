package tachymeter_test

import (
	"testing"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func BenchmarkAddTime(b *testing.B) {
	b.StopTimer()

	ta := tachymeter.New(&tachymeter.Config{Size: b.N})
	d := time.Millisecond

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.AddTime(d)
	}
}

func BenchmarkAddTimeSampling(b *testing.B) {
	b.StopTimer()

	ta := tachymeter.New(&tachymeter.Config{Size: 100})
	d := time.Millisecond

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ta.AddTime(d)
	}
}

func TestReset(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.Reset()

	if ta.Count != 0 {
		t.Fail()
	}
}

func TestAddTime(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.AddTime(time.Millisecond)

	if ta.Times[0] != time.Millisecond {
		t.Fail()
	}
}

func TestSetWallTime(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.SetWallTime(time.Millisecond)

	if ta.WallTime != time.Millisecond {
		t.Fail()
	}
}
