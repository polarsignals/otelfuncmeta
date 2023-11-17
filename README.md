# otelfuncmeta

This project provides an OpenTelemetry span processor for Go, that attaches source code metadata to every span created. This includes the function name, ffilename, line number and program-counter (pc).

## Usage

The span processor can either be passed at creation time of the tracer provider.

```go
import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "github.com/polarsignals/otelfuncmeta"
)

func newTracerProvider(t *testing.T, e *testExporter) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(otelfuncmeta.NewSpanProcessor()),
	)
	return tp
}
```

Or get registered after creating the tracer provider using `RegisterSpanProcessor`:

```go
import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "github.com/polarsignals/otelfuncmeta"
)

func newTracerProvider(t *testing.T, e *testExporter) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider()
    tp.RegisterSpanProcessor(otelfuncmeta.NewSpanProcessor())
	return tp
}
```

## Benchmarks

On an Apple Silicon M1 Pro, retrieving the function metadata takes about 190 nanoseconds.

```
$ go test -bench=. -count=5
goos: darwin
goarch: arm64
pkg: github.com/polarsignals/otelfuncmeta
BenchmarkFuncMetadata-10         6355422               190.1 ns/op             0 B/op          0 allocs/op
BenchmarkFuncMetadata-10         6274279               189.7 ns/op             0 B/op          0 allocs/op
BenchmarkFuncMetadata-10         6323715               189.5 ns/op             0 B/op          0 allocs/op
BenchmarkFuncMetadata-10         6275026               190.1 ns/op             0 B/op          0 allocs/op
BenchmarkFuncMetadata-10         6306594               189.8 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/polarsignals/otelfuncmeta    7.314s
```
