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
		c.AddCount(1)
	}

	results := c.Json()
	fmt.Println(results)

	c.Dump()
}
