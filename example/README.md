### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example`

### Run
```
$ $GOPATH/bin/example
{"Time":{"Cumulative":"705.24222ms","Avg":"14.104844ms","P50":"13.073198ms","P75":"21.358238ms","P95":"28.289403ms","P99":"30.544326ms","P999":"30.544326ms","Long5p":"29.843555ms","Short5p":"356.145µs","Max":"30.544326ms","Min":"2.455µs","Range":"30.541871ms"},"Rate":{"Second":70.89762720104873},"Samples":50,"Count":100}

50 samples of 100 events
Cumulative:     705.24222ms
Avg.:           14.104844ms
p50:            13.073198ms
p75:            21.358238ms
p95:            28.289403ms
p99:            30.544326ms
p999:           30.544326ms
Long 5%:        29.843555ms
Short 5%:       356.145µs
Max:            30.544326ms
Min:            2.455µs
Range:          30.541871ms
Rate/sec.:      70.90
```