# tachymeter
A simple latency summary library for Go

```go
import "github.com/jamiealquiza/tachymeter"

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 50})

	for i := 0; i < 100; i++ {
		start := time.Now()
		doSomeWork()
		c.AddTime(time.Since(start))
		c.AddCount(1)
	}

	c.Dump()
}
```

```
50 samples of 100 events
Total:			746.317344ms
Avg.:			14.926346ms
Median: 		13.902036ms
95%ile:			30.33816ms
Longest 5%:		31.541162ms
Shortest 5%:	466.619µs
Max:			32.146993ms
Min:			3.444µs
Rate/sec.:		67.00
```
