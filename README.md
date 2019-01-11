# queue

This is a FIFO queue for Go based on `container/list`.

### Features
- first in first out -- of course ;)
- goroutine safe
- both: blocking and non-blocking reading
- `Close` to release blocking
- all data is still read on `Close`

### Benchmarks

```
BenchmarkQueue_Push-12                	10000000	       122 ns/op	      48 B/op	       1 allocs/op
BenchmarkQueue_Pop-12                 	30000000	      42.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkQueue_PopBlocking-12         	30000000	      42.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkQueue_Push_PopBlocking-12    	10000000	       129 ns/op	      48 B/op	       1 allocs/op
```
