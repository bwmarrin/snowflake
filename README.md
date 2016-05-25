# flake
Very simple snowflake generator

Finally, hit my goal of 244ns/op.

But, interesting comparison of different methods.


```
time go test -bench=.
testing: warning: no tests to run
PASS
BenchmarkGenerateChan-12                 3000000               503 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateChanParallel-12         2000000               743 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateNoSleep-12              5000000               244 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateNoSleepLock-12          5000000               283 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateNoSleepLockParallel-12  5000000               348 ns/op               0 B/op          0 allocs/op
BenchmarkGenerate-12                     5000000               283 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateLocks-12                5000000               293 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateLocksParallel-12        5000000               368 ns/op               0 B/op          0 allocs/op
ok      _/home/bruce/flake      15.291s
go test -bench=.  16.88s user 7.37s system 151% cpu 15.981 total
```
