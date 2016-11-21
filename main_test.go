package tachymeter_test

import (
	"time"
	"testing"

	"github.com/jamiealquiza/tachymeter"
)

func TestReset(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.AddTime(time.Second)
	ta.AddCount(3)
	ta.Reset()

	if ta.TimesPosition != 0 {
		t.Fail()
	}
	if ta.TimesUsed != 0 {
		t.Fail()
	}
	if ta.Count != 0 {
		t.Fail()
	}
	if ta.TimeTotal != 0 {
		t.Fail()
	}
}

func TestAddTime(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.AddTime(time.Millisecond)

	if ta.TimeTotal != time.Millisecond {
		t.Fail()
	}
}

func TestAddCount(t *testing.T) {
	ta := tachymeter.New(&tachymeter.Config{Size: 3})

	ta.AddCount(3)

	if ta.Count != 3 {
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