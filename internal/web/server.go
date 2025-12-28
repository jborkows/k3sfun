package web

import (
	"net/http"
	"time"

	"shopping/internal/domain/admin"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/config"
	"shopping/internal/infrastructure/oidc"
)

// Timeout constants for consistent request handling across all handlers.
const (
	// DefaultHandlerTimeout is the default timeout for HTTP handlers.
	DefaultHandlerTimeout = 5 * time.Second
	// ShortHandlerTimeout is used for quick operations.
	ShortHandlerTimeout = 3 * time.Second
	// LongHandlerTimeout is used for potentially slower operations.
	LongHandlerTimeout = 10 * time.Second
)

type Server struct {
	cfg      config.Config
	auth     oidc.Authenticator
	products productsComponent
	shopping shoppingComponent
	admin    adminComponent
	events   *eventHub
	staticV  string
	units    []products.Unit // cached units loaded at startup
}

func NewServer(cfg config.Config, qry products.Queries, svc *products.Service, shoppingSvc *shoppinglist.Service, adminMaintenance admin.Maintenance, authenticator oidc.Authenticator, staticVersion string, units []products.Unit) *Server {
	return &Server{
		cfg:  cfg,
		auth: authenticator,
		products: productsComponent{
			qry: qry,
			svc: svc,
		},
		shopping: shoppingComponent{
			svc: shoppingSvc,
		},
		admin: adminComponent{
			maintenance: adminMaintenance,
		},
		events:  newEventHub(),
		staticV: staticVersion,
		units:   units,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.Handle("GET /login", http.HandlerFunc(s.auth.HandleLogin))
	mux.Handle("GET /oauth2/callback", http.HandlerFunc(s.auth.HandleCallback))
	mux.Handle("POST /logout", http.HandlerFunc(s.auth.HandleLogout))

	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/products", http.StatusFound) }))

	wrap := s.auth.Middleware

	s.registerProductRoutes(mux, wrap)
	s.registerAdminRoutes(mux, wrap)

	// Apply middleware chain: OpenTelemetry tracing -> Request logging -> Router
	return ChainMiddleware(mux,
		OtelMiddleware("shopping"),
		LoggingMiddleware,
	)
}
