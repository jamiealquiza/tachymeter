package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})
	tl := tachymeter.Timeline{}

	// Run 3 iterations of a loop that we're
	// interesting in summarizing with tachymeter.
	for iter := 0; iter < 3; iter++ {
		fmt.Printf("Running iteration %d\n", iter)
		// Capture timing data from the loop.
		for i := 0; i < 100; i++ {
			start := time.Now()
			time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
			c.AddTime(time.Since(start))
		}

		// Add each loop tachymeter
		// to the event timeline.
		tl.AddEvent(c.Calc())
		c.Reset()
	}

	// Write out an HTML page with the
	// histogram for all iterations.
	err := tl.WriteHTML(".")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Results written")
}
