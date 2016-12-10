### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example`

### Run
```
$ $GOPATH/bin/example
{"Time":{"Cumulative":"649.553805ms","Avg":"12.991076ms","P50":"12.062832ms","P75":"18.256032ms","P95":"27.175521ms","P99":"28.182451ms","P999":"28.182451ms","Long5p":"28.170356ms","Short5p":"416.962µs","Max":"28.182451ms","Min":"2.24µs"},"Rate":{"Second":76.97591733759454},"Samples":50,"Count":100}

50 samples of 100 events
Cumulative:	649.553805ms
Avg.:		12.991076ms
p50: 		12.062832ms
p75:		18.256032ms
p95:		27.175521ms
p99:		28.182451ms
p999:		28.182451ms
Long 5%:	28.170356ms
Short 5%:	416.962µs
Max:		28.182451ms
Min:		2.24µs
Rate/sec.:	76.98
```