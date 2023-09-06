package cloudtrace

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dustedcodes/blog/internal/scrub"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var Middleware = func(next http.Handler) http.Handler {
	return otelhttp.NewHandler(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())

				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.method"),
					Value: attribute.StringValue(r.Method),
				})
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.url"),
					Value: attribute.StringValue(r.URL.String()),
				})
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.protocol"),
					Value: attribute.StringValue(r.Proto),
				})
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.remote_ip"),
					Value: attribute.StringValue(r.RemoteAddr),
				})
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.user_agent"),
					Value: attribute.StringValue(r.UserAgent()),
				})
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.referer"),
					Value: attribute.StringValue(r.Referer()),
				})
				headers := []string{}
				for k, v := range r.Header {
					value := strings.Join(v, ",")
					if k == "Authorization" || k == "Cookie" {
						value = scrub.Middle(value)
					}
					headers = append(headers, fmt.Sprintf("%s: %s", k, value))
				}
				span.SetAttributes(attribute.KeyValue{
					Key:   attribute.Key("http.headers"),
					Value: attribute.StringSliceValue(headers),
				})

				next.ServeHTTP(w, r)
			},
		),
		"http-request",
		otelhttp.WithSpanNameFormatter(
			func(_ string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			}),
	)
}
