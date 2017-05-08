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

	// Print JSON format to console.
	resultsj := c.Json()
	fmt.Printf("%s\n\n", resultsj)

	// Print pre-formatted console output.
	c.Calc().Dump()
}