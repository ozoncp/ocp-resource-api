package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api "github.com/ozoncp/ocp-resource-api/internal/api"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	desc "github.com/ozoncp/ocp-resource-api/pkg/ocp-resource-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

const (
	grpcPort = ":7072"
	httpPort = ":7070"
)

var (
	grpcEndpoint = flag.String("grpc-server-endpoint", "0.0.0.0"+grpcPort, "gRPC server endpoint")
	httpEndpoint = flag.String("http-server-endpoint", "0.0.0.0"+httpPort, "HTTP server endpoint")
)

func runGrpcServer() error {
	l, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal().Msgf("failed to close: %v", err)
		}
	}(l)

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal().Msgf("DB_URL environment variable should be defined")
	}

	grpcServer := grpc.NewServer()
	resourceRepo := repo.NewRepoPostgreSQL(db_url)
	resourceApi, err := api.NewOcpResourceApi(&resourceRepo)
	if err != nil {
		panic(err)
	}
	desc.RegisterOcpResourceApiServer(grpcServer, resourceApi)
	log.Info().Msgf("Started grpc server on %s", *grpcEndpoint)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}

	return nil
}

func runHttpServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpResourceApiHandlerFromEndpoint(ctx, mux, *grpcEndpoint, opts)
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Start to listen http server on %s", *httpEndpoint)
	err = http.ListenAndServe(*httpEndpoint, mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	go runHttpServer()

	if err := runGrpcServer(); err != nil {
		log.Fatal().Err(err)
	}
}
