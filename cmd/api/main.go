package main

import (
	"log"
	"net"
	"net/http"

	"pvz-service/api/grpc/pvz/pvz_v1"
	"pvz-service/internal/config"
	"pvz-service/internal/db"
	"pvz-service/internal/metrics"
	"pvz-service/internal/middleware"
	"pvz-service/internal/transport/grpc/pvz"
	"pvz-service/internal/transport/handlers"
	"pvz-service/pkg/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startHTTPServer(address string) {
	router := handlers.RegisterRoutes()

	log.Printf("HTTP-сервер запущен на %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}

func startGRPCServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GrpcMetrics))

	pvz_v1.RegisterPVZServiceServer(grpcServer, pvz.NewPVZServer())
	reflection.Register(grpcServer)

	log.Printf("gRPC-сервер запущен на %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка gRPC-сервера: %v", err)
	}
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Метрики доступны на :9000/metrics")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера метрик: %v", err)
	}
}

func main() {
	logger.Init()
	metrics.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	dbPool, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer dbPool.Close()

	go startHTTPServer(cfg.HTTP.ListenAddress())
	go startGRPCServer(cfg.GRPC.ListenAddress())
	startMetricsServer()
}
