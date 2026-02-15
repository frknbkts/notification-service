package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/frknbkts/notification-service/internal/middleware"
	"github.com/frknbkts/notification-service/internal/repository"
	"github.com/frknbkts/notification-service/internal/service"
	"github.com/frknbkts/notification-service/pkg/pb"
)

func main() {

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":50051"
	}

	cbConnStr := os.Getenv("COUCHBASE_CONNECTION_STRING")
	if cbConnStr == "" {
		cbConnStr = "couchbase://localhost"
	}

	cbUser := os.Getenv("COUCHBASE_USERNAME")
	if cbUser == "" {
		cbUser = "admin"
	}

	cbPass := os.Getenv("COUCHBASE_PASSWORD")
	if cbPass == "" {
		cbPass = "password123"
	}

	cbBucket := os.Getenv("COUCHBASE_BUCKET_NAME")
	if cbBucket == "" {
		cbBucket = "notification"
	}

	fmt.Println("Couchbase baglantisi kuruluyor...")
	repo, err := repository.NewCouchbaseRepository(cbConnStr, cbUser, cbPass, cbBucket)
	if err != nil {
		log.Fatalf("Veritabani baglanti hatasi: %v", err)
	}
	fmt.Println("Couchbase baglantisi basarili.")

	notificationService := service.NewNotificationService(repo)
	rateLimiter := middleware.NewRateLimiterInterceptor(rate.Limit(100), 10)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggerInterceptor(),
			rateLimiter.Unary(),
			middleware.MetricsInterceptor,
		),
		grpc.ChainStreamInterceptor(
			rateLimiter.Stream(),
		),
	)

	pb.RegisterNotificationServiceServer(grpcServer, notificationService)
	reflection.Register(grpcServer)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Metrics server :9090 portunda calisiyor...")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Printf("Metrics sunucusu hatasi: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("gRPC Sunucusu %s portunda calisiyor...\n", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Sunucu hatasi: %v", err)
		}
	}()

	<-stopChan
	fmt.Println("\nKapanma sinyali alindi...")
	grpcServer.GracefulStop()
	fmt.Println("Sunucu kapatildi.")
}
