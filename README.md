# flake
Very simple snowflake generator

but..

Why is the mutex'd version faster?

```
go test -bench=.
testing: warning: no tests to run
PASS
BenchmarkGenerate-12                     3000000               523 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateLocks-12                5000000               290 ns/op               0 B/op          0 allocs/op
BenchmarkGenerateLocksParallel-12        3000000               540 ns/op               0 B/op          0 allocs/op
```
