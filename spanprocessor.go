package otelfuncmeta

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type SpanProcessor struct {
	callerDepth int
}

func NewSpanProcessor() *SpanProcessor {
	return &SpanProcessor{
		callerDepth: 4,
	}
}

func NewSpanProcessorWithCallerDepth(depth int) *SpanProcessor {
	return &SpanProcessor{
		callerDepth: depth,
	}
}

// OnStart is called when a span is started.
func (s *SpanProcessor) OnStart(parent context.Context, span sdktrace.ReadWriteSpan) {
	pc, funcName, file, line := funcMetadata(s.callerDepth)
	span.SetAttributes(
		attribute.String("function.pc", fmt.Sprintf("0x%x", pc)),
		attribute.String("function.name", funcName),
		attribute.String("function.file", file),
		attribute.Int("function.line", line),
	)
}

func funcMetadata(depth int) (uintptr, string, string, int) {
	var pcbuf [1]uintptr
	runtime.Callers(depth, pcbuf[:])
	pc := pcbuf[0]
	f := runtime.FuncForPC(pc)
	funcName := f.Name()
	file, line := f.FileLine(pc)
	return pc, funcName, file, line
}

func (s *SpanProcessor) OnEnd(_ sdktrace.ReadOnlySpan) {}

func (s *SpanProcessor) Shutdown(_ context.Context) error { return nil }

func (s *SpanProcessor) ForceFlush(_ context.Context) error { return nil }
