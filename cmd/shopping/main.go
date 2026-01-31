package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"shopping/internal/domain/admin"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/config"
	"shopping/internal/infrastructure/logging"
	"shopping/internal/infrastructure/oidc"
	"shopping/internal/infrastructure/persistence/sqlite"
	"shopping/internal/migrator"
	"shopping/internal/web"
)

var buildVersion = ""

func main() {
	logger := logging.New()
	slog.SetDefault(logger)

	staticVersion := strings.TrimSpace(buildVersion)
	if staticVersion == "" {
		staticVersion = strconv.FormatInt(time.Now().Unix(), 10)
	}

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Log database path (without sensitive info)
	dsnForLog := cfg.DBDSN
	if idx := strings.Index(dsnForLog, "?"); idx != -1 {
		dsnForLog = dsnForLog[:idx]
	}
	slog.Info("Connecting to database", "dsn", dsnForLog)

	conn, err := sqlite.Open(cfg.DBDSN)
	if err != nil {
		slog.Error("Failed to open database", "error", err, "dsn", dsnForLog)
		log.Fatalf("db: %v", err)
	}
	defer conn.Close()
	slog.Info("Database connection established")

	if err := migrator.Up(conn); err != nil {
		slog.Error("Migrations failed", "error", err)
		log.Fatalf("migrate: %v", err)
	}
	slog.Info("Database migrations completed")

	repo := sqlite.NewRepo(conn)
	var productsQueries products.Queries = repo
	productsService := products.NewService(repo)
	var adminMaintenance admin.Maintenance = repo
	var shoppingRepo shoppinglist.Repository = repo
	shoppingService := shoppinglist.NewService(shoppingRepo, productsService)

	// Load units once at startup
	units, err := productsQueries.ListUnits(context.Background())
	if err != nil {
		log.Fatalf("load units: %v", err)
	}

	authenticator, err := oidc.New(cfg)
	if err != nil {
		log.Fatalf("auth: %v", err)
	}

	srv := web.NewServer(cfg, productsQueries, productsService, shoppingService, adminMaintenance, authenticator, staticVersion, units)

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}
