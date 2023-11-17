package otelfuncmeta

import (
	"context"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func basicTracerProvider(t *testing.T, e *testExporter) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(NewSpanProcessor()),
		sdktrace.WithSyncer(e),
	)
	return tp
}

type testExporter struct {
	spans []sdktrace.ReadOnlySpan
}

func (e *testExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *testExporter) Shutdown(_ context.Context) error {
	return nil
}

func TestSpanProcessor(t *testing.T) {
	e := &testExporter{}
	tp := basicTracerProvider(t, e)
	tr := tp.Tracer("NilExporter")

	_, span := tr.Start(context.Background(), "foo")
	span.End()

	if len(e.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(e.spans))
	}

	attrs := e.spans[0].Attributes()
	if len(attrs) != 4 {
		t.Fatalf("expected 4 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "function.pc" {
		t.Fatalf("expected function.pc, got %s", attrs[0].Key)
	}
	// Intentionally not testing the value of the program counter as it is
	// likely to change across Go versions.
	if attrs[1].Key != "function.name" {
		t.Fatalf("expected function.name, got %s", attrs[1].Key)
	}
	if attrs[1].Value.AsString() != "github.com/polarsignals/otelfuncmeta.TestSpanProcessor" {
		t.Fatalf("expected TestSpanProcessor, got %s", attrs[1].Value.AsString())
	}
	if attrs[2].Key != "function.file" {
		t.Fatalf("expected function.file, got %s", attrs[2].Key)
	}
	if !strings.HasSuffix(attrs[2].Value.AsString(), "spanprocessor_test.go") {
		t.Fatalf("expected to end with spanprocessor_test.go, got %s", attrs[2].Value.AsString())
	}
	if attrs[3].Key != "function.line" {
		t.Fatalf("expected function.line, got %s", attrs[3].Key)
	}
	if attrs[3].Value.AsInt64() != 39 {
		t.Fatalf("expected 38, got %d", attrs[3].Value.AsInt64())
	}
}

func BenchmarkFuncMetadata(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = funcMetadata(1)
	}
}
