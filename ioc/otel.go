package ioc

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

func InitOtel(l logx.Logger) func(ctx context.Context) {
	res, err := newResource("ink-flow", "v0.0.1")
	if err != nil {
		panic(err)
	}
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		er := tp.Shutdown(ctx)
		if er != nil {
			l.Error("shutdown otel error", logx.Error(er))
		} else {
			l.Info("shutdown otel success")
		}
	}
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion)),
	)
}
func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	// TODO 地址后续需要改成配置项
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exporter,
		trace.WithExportTimeout(time.Second)),
		trace.WithResource(res))
	return tracerProvider, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{})
}
