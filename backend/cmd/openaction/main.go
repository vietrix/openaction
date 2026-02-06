package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openaction/internal/api"
	"openaction/internal/auth"
	"openaction/internal/blob"
	"openaction/internal/config"
	"openaction/internal/db"
	"openaction/internal/pool"
	"openaction/internal/ui"
	"openaction/pkg/poolpb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		log.Fatalf("data dir error: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := db.Migrate(ctx, database, filepath.Clean("../backend/migrations")); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	authService := &auth.Service{
		DB:         database,
		SessionTTL: cfg.SessionTTL,
		TokenTTL:   cfg.TokenTTL,
		CSRFFlag:   cfg.CSRFEnabled,
	}
	if err := authService.EnsureAdmin(ctx, cfg.AdminEmail, cfg.AdminPass); err != nil {
		log.Fatalf("ensure admin error: %v", err)
	}

	blobStore := blob.New(cfg.DataDir)

	apiServer := &api.Server{
		DB:         database,
		Auth:       authService,
		Blob:       blobStore,
		DataDir:    cfg.DataDir,
		SecureOnly: cfg.TLSCertPath != "" && cfg.TLSKeyPath != "",
	}

	router := chi.NewRouter()
	router.Mount("/", apiServer.Router())
	if cfg.ServeUI {
		ui.Attach(router, cfg.UITargetDir)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if cfg.TLSCertPath != "" && cfg.TLSKeyPath != "" {
			log.Printf("HTTP server (TLS) listening on %s", httpServer.Addr)
			if err := httpServer.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("http server error: %v", err)
			}
			return
		}
		log.Printf("HTTP server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	grpcServer, grpcListener, err := startGRPC(cfg)
	if err != nil {
		log.Fatalf("grpc error: %v", err)
	}
	defer grpcListener.Close()
	defer grpcServer.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownCtx)
}

func startGRPC(cfg *config.Config) (*grpc.Server, net.Listener, error) {
	if cfg.TLSCertPath == "" || cfg.TLSKeyPath == "" || cfg.CACertPath == "" {
		return grpc.NewServer(), dummyListener{}, nil
	}

	cert, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
	if err != nil {
		return nil, nil, err
	}

	caCertPEM, err := os.ReadFile(cfg.CACertPath)
	if err != nil {
		return nil, nil, err
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertPEM) {
		return nil, nil, fmt.Errorf("failed to parse ca cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
		MinVersion:   tls.VersionTLS13,
	}

	creds := credentials.NewTLS(tlsConfig)
	server := grpc.NewServer(grpc.Creds(creds))
	poolpb.RegisterPoolServiceServer(server, &pool.Server{})

	listener, err := net.Listen("tcp", cfg.PoolGRPCAddr)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		log.Printf("gRPC server listening on %s", cfg.PoolGRPCAddr)
		if err := server.Serve(listener); err != nil {
			log.Printf("grpc serve error: %v", err)
		}
	}()

	return server, listener, nil
}

type dummyListener struct{}

func (d dummyListener) Accept() (net.Conn, error) { return nil, fmt.Errorf("grpc disabled") }
func (d dummyListener) Close() error              { return nil }
func (d dummyListener) Addr() net.Addr            { return dummyAddr("") }

type dummyAddr string

func (d dummyAddr) Network() string { return "tcp" }
func (d dummyAddr) String() string  { return string(d) }
