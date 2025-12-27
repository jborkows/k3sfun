package web

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// responseWriter wraps http.ResponseWriter to capture status code.
// It also implements http.Flusher to support SSE streaming.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Flush implements http.Flusher by delegating to the underlying writer.
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// LoggingMiddleware logs HTTP requests with method, path, status, and duration.
// It skips wrapping for SSE endpoints to preserve http.Flusher support.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip wrapping for SSE endpoints to preserve Flusher interface
		if isSSEPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Skip logging for health checks and static assets to reduce noise
		if r.URL.Path == "/healthz" || isStaticPath(r.URL.Path) {
			return
		}

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", duration,
			"remote_addr", r.RemoteAddr,
		)
	})
}

// isStaticPath returns true if the path is for static assets.
func isStaticPath(path string) bool {
	return len(path) > 8 && path[:8] == "/static/"
}

// isSSEPath returns true if the path is for Server-Sent Events.
func isSSEPath(path string) bool {
	return strings.HasPrefix(path, "/events")
}

// OtelMiddleware wraps the handler with OpenTelemetry instrumentation.
// It skips SSE endpoints to preserve http.Flusher support for streaming.
func OtelMiddleware(serviceName string) func(http.Handler) http.Handler {
	otelHandler := otelhttp.NewMiddleware(serviceName,
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip otel wrapping for SSE endpoints to preserve Flusher interface
			if isSSEPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			otelHandler(next).ServeHTTP(w, r)
		})
	}
}

// ChainMiddleware chains multiple middleware together.
func ChainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
