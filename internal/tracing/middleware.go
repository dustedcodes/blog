package tracing

import (
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func scrub(s string) string {
	const mask = "********"

	l := len(s)
	if l <= 8 {
		return mask
	}
	if l > 8 && l <= 10 {
		return s[:2] + mask + s[l-2:]
	}
	if l > 10 && l < 20 {
		return s[:3] + mask + s[l-3:]
	}
	if l >= 20 && l < 30 {
		return s[:5] + mask + s[l-5:]
	}
	return s[:10] + mask + s[l-10:]
}

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
						value = scrub(value)
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
