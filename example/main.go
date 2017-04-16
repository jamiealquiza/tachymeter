package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})

	for i := 0; i < 100; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		c.AddTime(time.Since(start))
	}

	resultsj := c.Json()
	fmt.Printf("%s\n\n", resultsj)

	results := c.Calc()
	// Print the pre-formatted console output.
	results.Dump()
	// Create an HTML graph of the event histogram.
	//results.DumpHistogramGraph()
	w := tachymeter.Timeline{}
	w.AddEvent(results)
	w.WriteHtml()
}
