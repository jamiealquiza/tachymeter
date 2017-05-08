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

	for iter := 0; iter <3; iter++ {
		fmt.Printf("Running iteration %d\n", iter)
		for i := 0; i < 100; i++ {
			start := time.Now()
			time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
			c.AddTime(time.Since(start))
		}

		tl.AddEvent(c.Calc())
		c.Reset()
	}


	// Create an HTML graph of the event histogram.
	err := tl.WriteHtml(".")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Results written")
}
