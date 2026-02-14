package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/frknbkts/notification-service/internal/repository"
	"github.com/frknbkts/notification-service/internal/service"
	"github.com/frknbkts/notification-service/pkg/pb"
)

func main() {
	
	port := ":50051"
	cbConnStr := "couchbase://localhost"
	cbUser := "admin"
	cbPass := "password123"
	cbBucket := "notification"

	fmt.Println("Couchbase baglantisi kuruluyor...")
	repo, err := repository.NewCouchbaseRepository(cbConnStr, cbUser, cbPass, cbBucket)
	if err != nil {
		log.Fatalf("Veritabani baglanti hatasi: %v", err)
	}
	fmt.Println("Couchbase baglantisi basarili.")

	notificationService := service.NewNotificationService(repo)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterNotificationServiceServer(grpcServer, notificationService)

	reflection.Register(grpcServer)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("gRPC Sunucusu %s portunda calisiyor...\n", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Sunucu hatasi: %v", err)
		}
	}()

	<-stopChan
	fmt.Println("\nKapanma sinyali alindi. Graceful shutdown baslatiliyor...")

	grpcServer.GracefulStop()
	fmt.Println("Sunucu basariyla kapatildi.")
}