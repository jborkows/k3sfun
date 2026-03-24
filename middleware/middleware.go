// Package middleware provides HTTP middleware for logging, OpenTelemetry tracing,
// and request handling utilities.
package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

type Config struct {
	Logger          *slog.Logger
	PublicPaths     []string
	PublicPrefixes  []string
	SSEPathPrefixes []string
}

type Option func(*Config)

func WithLogger(l *slog.Logger) Option {
	return func(c *Config) {
		c.Logger = l
	}
}

func PublicPath(path string) Option {
	return func(c *Config) {
		c.PublicPaths = append(c.PublicPaths, path)
	}
}

func PublicPathPrefix(prefix string) Option {
	return func(c *Config) {
		c.PublicPrefixes = append(c.PublicPrefixes, prefix)
	}
}

func WithSSEPathPrefixes(prefixes ...string) Option {
	return func(c *Config) {
		c.SSEPathPrefixes = prefixes
	}
}

func defaultConfig() *Config {
	return &Config{
		Logger:          slog.Default(),
		PublicPaths:     []string{"/healthz", "/favicon.ico"},
		PublicPrefixes:  []string{"/static/"},
		SSEPathPrefixes: []string{"/events"},
	}
}

func (c *Config) isPublicPath(path string) bool {
	for _, p := range c.PublicPaths {
		if path == p {
			return true
		}
	}
	for _, prefix := range c.PublicPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func (c *Config) isSSEPath(path string) bool {
	for _, prefix := range c.SSEPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func LoggingMiddleware(opts ...Option) func(http.Handler) http.Handler {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.isSSEPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			if cfg.isPublicPath(r.URL.Path) {
				return
			}

			cfg.Logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", duration,
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

func OtelMiddleware(serviceName string, opts ...Option) func(http.Handler) http.Handler {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	otelHandler := otelhttp.NewMiddleware(serviceName,
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.isSSEPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			otelHandler(next).ServeHTTP(w, r)
		})
	}
}

func ChainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
