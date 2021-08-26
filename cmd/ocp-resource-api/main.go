package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	api "github.com/ozoncp/ocp-resource-api/internal/api"
	"github.com/ozoncp/ocp-resource-api/internal/metrics"
	"github.com/ozoncp/ocp-resource-api/internal/producer"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	"github.com/ozoncp/ocp-resource-api/internal/tracer"
	desc "github.com/ozoncp/ocp-resource-api/pkg/ocp-resource-api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

const (
	grpcPort   = ":7072"
	httpPort   = ":7070"
	metricPort = ":9100"
)

var (
	grpcEndpoint   = flag.String("grpc-server-endpoint", "0.0.0.0"+grpcPort, "gRPC server endpoint")
	httpEndpoint   = flag.String("http-server-endpoint", "0.0.0.0"+httpPort, "HTTP server endpoint")
	metricEndpoint = flag.String("metric-endpoint", "0.0.0.0"+metricPort, "Metric server endpoint")
	kafka_url      = flag.String("kafka_url", "127.0.0.1:9094", "Connection URL for DB")
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

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal().Msgf("DB_URL environment variable should be defined")
	}

	grpcServer := grpc.NewServer()
	resourceRepo := repo.NewRepoPostgreSQL(dbUrl)
	jaegerTracer, jaegerCloser, err := tracer.CreateTracer()
	if err != nil {
		log.Fatal().Msgf("Issue during tracer initialization: %v", err)
	}
	defer func() {
		closer := *jaegerCloser
		err := closer.Close()
		if err != nil {
			log.Err(err)
		}
	}()
	opentracing.SetGlobalTracer(*jaegerTracer)
	prod, err := producer.NewProducer([]string{*kafka_url}, "ocp-resource-api")
	if err != nil {
		log.Fatal().Msgf("Could not initialize kafka with %v brokers", *kafka_url)
	}

	resourceApi, err := api.NewOcpResourceApi(resourceRepo, prod)
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
	metricsSrv := runMetricsServer()
	metrics.RegisterMetrics()
	defer func(metricsSrv *http.Server) {
		err := metricsSrv.Close()
		if err != nil {
			log.Err(err)
		}
	}(metricsSrv)

	go runHttpServer()

	if err := runGrpcServer(); err != nil {
		log.Fatal().Err(err)
	}
}

func runMetricsServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("prom", promhttp.Handler())
	metricsSrv := &http.Server{
		Addr:    *metricEndpoint,
		Handler: mux,
	}
	return metricsSrv
}
