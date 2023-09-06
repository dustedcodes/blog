package cloudtrace

import (
	"context"
	"fmt"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	tracerName = "github.com/dustedcodes/msgdrop/internal/thirdparty/cloudtrace"
)

type Options struct {
	InitForGCP      bool
	GCPProjectID    string
	ServiceName     string
	ServiceVersion  string
	EnvironmentName string
}

// Resources used:
// ---
// https://opentelemetry.io/docs/instrumentation/go/getting-started/
// https://github.com/GoogleCloudPlatform/opentelemetry-operations-go
// https://github.com/open-telemetry/opentelemetry-go-contrib/

// OpenTelemetry in GCP
// ---
// Currently the Google Cloud Go libraries don't support OpenTelemetry out of the box.
// They were written with OTel's predecessor, OpenCensus, in mind.
// The OpenTelemetry project provides a bridge between the two:
// https://github.com/open-telemetry/opentelemetry-go/tree/main/bridge/opencensus

func getOTLPExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error creating OTLP trace exporter: %w", err)
	}
	return exporter, nil
}

func getGoogleCloudExporter(googleProjectID string) (sdktrace.SpanExporter, error) {
	exporter, err := gcpexporter.New(gcpexporter.WithProjectID(googleProjectID))
	if err != nil {
		return nil, fmt.Errorf("error creating GCP trace exporter: %w", err)
	}
	return exporter, nil
}

func getExporter(ctx context.Context, googleProjectID string, initForGCP bool) (sdktrace.SpanExporter, error) {
	if initForGCP {
		return getGoogleCloudExporter(googleProjectID)
	}
	return getOTLPExporter(ctx)
}

func getSampler(initForGCP bool) sdktrace.Sampler {
	if !initForGCP {
		return sdktrace.AlwaysSample()
	}
	return sdktrace.ParentBased(
		sdktrace.AlwaysSample(),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()))
}

func installPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

func initTracer(
	ctx context.Context,
	opts *Options,
) (
	*sdktrace.TracerProvider,
	error,
) {
	exporter, err := getExporter(ctx, opts.GCPProjectID, opts.InitForGCP)
	if err != nil {
		return nil, fmt.Errorf("error creating trace exporter: %w", err)
	}
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithAttributes(
			semconv.TelemetrySDKName("opentelemetry"),
			semconv.TelemetrySDKVersion(sdk.Version()),
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceNameKey.String(opts.ServiceName),
			semconv.ServiceVersionKey.String(opts.ServiceVersion),
			attribute.String("environment", opts.EnvironmentName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating tracing resource: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(getSampler(opts.InitForGCP)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

func Init(
	ctx context.Context,
	opts *Options,
) (
	func(),
	error,
) {
	installPropagators()
	tp, err := initTracer(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error initialising tracer: %w", err)
	}
	otel.SetTracerProvider(tp)
	shutdown := func() {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down trace provider: %+v", err)
		}
	}
	tracer := otel.Tracer(tracerName)
	octrace.DefaultTracer = opencensus.NewTracer(tracer)
	return shutdown, nil
}
