# localstorage

The package provides `key=value` storage with zero memory allocation.

## Benchmark

#### GOMAXPROCS=1

```sh
$ GOMAXPROCS=1 go test -bench=LocalStorage -benchmem -benchtime=10s
goos: darwin
goarch: amd64
pkg: github.com/serge64/localstorage
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLocalStorage_AsyncAll  45216859               260.1 ns/op             0 B/op          0 allocs/op
BenchmarkLocalStorage_Get       302990174               40.12 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Put       236891362               50.93 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Del       282445515               42.67 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Keys      777665688               15.34 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Values    760686696               15.74 ns/op            0 B/op          0 allocs/op
```

#### GOMAXPROCS=4

```sh
$ GOMAXPROCS=4 go test -bench=LocalStorage -benchmem -benchtime=10s
goos: darwin
goarch: amd64
pkg: github.com/serge64/localstorage
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLocalStorage_AsyncAll-4        23296315               484.1 ns/op             0 B/op          0 allocs/op
BenchmarkLocalStorage_Get-4             319952442               38.42 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Put-4             229826967               50.73 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Del-4             293215400               40.94 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Keys-4            767696656               15.36 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Values-4          764185610               15.38 ns/op            0 B/op          0 allocs/op
```

#### GOMAXPROCS=8

```sh
$ GOMAXPROCS=8 go test -bench=LocalStorage -benchmem -benchtime=10s
goos: darwin
goarch: amd64
pkg: github.com/serge64/localstorage
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLocalStorage_AsyncAll-8        19037619               626.4 ns/op             0 B/op          0 allocs/op
BenchmarkLocalStorage_Get-8             301694812               39.31 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Put-8             225441598               53.33 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Del-8             272035114               43.01 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Keys-8            740967951               15.78 ns/op            0 B/op          0 allocs/op
BenchmarkLocalStorage_Values-8          775615560               15.34 ns/op            0 B/op          0 allocs/op
```
