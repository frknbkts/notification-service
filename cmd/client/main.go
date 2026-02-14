package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/frknbkts/notification-service/pkg/pb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Baglanamadi: %v", err)
	}
	defer conn.Close()

	client := pb.NewNotificationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("--- Bildirim Gonderiliyor ---")
	sendRes, err := client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:   "user_123",
		SenderId: "user_999",
		Type:     pb.NotificationType_LIKE,
		Title:    "Yeni Begeni",
		Body:     "Ahmet fotrafini begendi",
	})
	if err != nil {
		log.Fatalf("Gonderim hatasi: %v", err)
	}
	fmt.Printf("✅ Bildirim Olusturuldu! ID: %s, Success: %v\n", sendRes.Id, sendRes.Success)

	fmt.Println("\n--- Bildirimler Listeleniyor ---")
	listRes, err := client.GetNotifications(ctx, &pb.GetNotificationsRequest{
		UserId: "user_123",
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		log.Fatalf("Listeleme hatasi: %v", err)
	}

	fmt.Printf("Toplam Bildirim Sayisi: %d\n", listRes.TotalCount)
	for _, n := range listRes.Notifications {
		fmt.Printf("- [%s] %s: %s (Okundu: %v)\n", n.Id, n.Title, n.Body, n.IsRead)
	}
}