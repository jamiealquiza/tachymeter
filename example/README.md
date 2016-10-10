### Install
 - `$ go get github.com/jamiealquiza/tachymeter`
 - `$ go install github.com/jamiealquiza/tachymeter/example`

### Run
```
$ $GOPATH/bin/example
{"Time":{"Total":"656.49102ms","Avg":"13.12982ms","Median":"12.150708ms","Long5p":"27.508286ms","Short5p":"438.305µs","Max":"28.153238ms","Min":"3.242µs"},"Rate":{"Second":76.16250409640028},"Samples":50,"Count":100}

50 samples of 100 events
Total:		656.49102ms
Avg.:		13.12982ms
Median: 	12.150708ms
95%ile:		26.466598ms
Longest 5%:	27.508286ms
Shortest 5%:	438.305µs
Max:		28.153238ms
Min:		3.242µs
Rate/sec.:	76.16
```