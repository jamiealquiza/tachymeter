### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example`

### Run
```
$ $GOPATH/bin/example
{"Time":{"Total":"1.463173043s","Avg":"29.26346ms","Median":"12.037232ms","Long5p":"27.402056ms","Short5p":"383.674µs","Max":"28.024171ms","Min":"391ns"},"Rate":{"Second":68.34461616034571},"Samples":50,"Count":100}

50 samples of 100 events
Total:			1.463173043s
Avg.:			29.26346ms
Median: 		12.037232ms
95%ile:			27.036779ms
Longest 5%:		27.402056ms
Shortest 5%:	383.674µs
Max:			28.024171ms
Min:			391ns
Rate/sec.:		68.34
```