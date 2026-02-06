package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openaction/pkg/poolpb"
)

func main() {
	addr := os.Getenv("OA_POOL_ADDR")
	if addr == "" {
		addr = "127.0.0.1:7443"
	}

	tlsConfig, err := loadTLS()
	if err != nil {
		log.Fatalf("tls error: %v", err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("grpc dial error: %v", err)
	}
	defer conn.Close()

	client := poolpb.NewPoolServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &poolpb.RegisterRequest{
		Info: &poolpb.PoolInfo{
			Id:      "pool-dev",
			Name:    "local-pool",
			Version: "0.1.0",
		},
	})
	if err != nil {
		log.Fatalf("register error: %v", err)
	}

	log.Printf("registered pool: %s", resp.AssignedId)
}

func loadTLS() (*tls.Config, error) {
	certPath := os.Getenv("OA_POOL_CERT")
	keyPath := os.Getenv("OA_POOL_KEY")
	caPath := os.Getenv("OA_POOL_CA")

	if certPath == "" || keyPath == "" || caPath == "" {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	caData, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caData) {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
