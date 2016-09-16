# tachymeter
A simple latency summary library for Go

```go
import "github.com/jamiealquiza/tachymeter"

func main() {
	c := tachymeter.New(&tachymeter.Config{Size: 500)

	for i := 0; i < 1000; i++ {
		start := time.Now()
		doSomeWork()
		c.AddTime(time.Since(start))
		c.AddCount(1)
	}

	c.Dump()
}
```

```
500 samples of 1000 events
Total:			45.466810987s
Avg.:			4.546681ms
95%ile:			2.612486ms
Longest 5%:	25.997116ms
Shortest 5%:	1.914673ms
Max:			19.555149795s
Min:			755.551µs
Rate/sec.:		10.99
```
