package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})

	// Measure events.
	for i := 0; i < 100; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		c.AddTime(time.Since(start))
	}

	// Calc output.
	results := c.Calc()

	// Print JSON format to console.
	fmt.Printf("%s\n\n", results.JSON())

	// Print pre-formatted console output.
	fmt.Printf("%s\n\n", results)

	// Print text histogram.
	fmt.Println(results.Histogram.String(25))
}
