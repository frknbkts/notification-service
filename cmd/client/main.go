package main

import (
	"context"
	"fmt"
	"io"
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

	streamCtx := context.Background()
	rpcCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("NOTIFICATION CLIENT BASLATILIYOR...")

	go func() {
		fmt.Println("[Stream] Canlı Akışa Bağlanılıyor...")
		stream, err := client.StreamNotifications(streamCtx, &pb.StreamNotificationsRequest{UserId: "user_123"})
		if err != nil {
			log.Printf("Stream hatasi: %v", err)
			return
		}

		for {
			notification, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Stream okuma hatasi: %v", err)
				return
			}
			fmt.Printf("\n>>> CANLI YAYIN: Yeni Bildirim Geldi! <<<\n")
			fmt.Printf("   Başlık: %s\n   İçerik: %s\n\n", notification.Title, notification.Body)
		}
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("--- 1. Yeni Bildirim Gonderiliyor (SendNotification) ---")
	sendRes, err := client.SendNotification(rpcCtx, &pb.SendNotificationRequest{
		UserId:   "user_123",
		SenderId: "system_bot",
		Type:     pb.NotificationType_SYSTEM,
		Title:    "Sistem Testi",
		Body:     "Bu bildirim Client testi sirasinda otomatik olusturuldu.",
	})
	if err != nil {
		log.Fatalf("Gonderim hatasi: %v", err)
	}
	fmt.Printf("Bildirim Kaydedildi! ID: %s\n", sendRes.Id)

	time.Sleep(2 * time.Second)

	fmt.Println("\n--- 2. Bildirimler Listeleniyor (GetNotifications) ---")
	listRes, err := client.GetNotifications(rpcCtx, &pb.GetNotificationsRequest{
		UserId: "user_123",
		Page:   1,
		Limit:  5,
	})
	if err != nil {
		log.Fatalf("Listeleme hatasi: %v", err)
	}

	var targetID string
	for _, n := range listRes.Notifications {
		status := "OKUNMADI"
		if n.IsRead {
			status = "OKUNDU"
		}

		prefix := "-"
		if n.Id == sendRes.Id {
			prefix = "->"
			targetID = n.Id
		}
		fmt.Printf("%s [%s] %s (%s)\n", prefix, n.Id, n.Title, status)
	}

	if targetID == "" {
		log.Fatalf("Hata: Az once ekledigimiz bildirim listede gorunmuyor!")
	}

	fmt.Printf("\n--- 3. Bildirim Okunuyor: %s (MarkAsRead) ---\n", targetID)
	readRes, err := client.MarkAsRead(rpcCtx, &pb.MarkAsReadRequest{
		UserId:         "user_123",
		NotificationId: targetID,
	})
	if err != nil {
		log.Fatalf("Okuma hatasi: %v", err)
	}
	fmt.Printf("Islem Basarili mi: %v\n", readRes.Success)

	fmt.Println("\n--- 4. Kontrol (GetUnreadCount & Verify) ---")

	countRes, _ := client.GetUnreadCount(rpcCtx, &pb.GetUnreadCountRequest{UserId: "user_123"})
	fmt.Printf("Kalan Okunmamis Sayisi: %d\n", countRes.Count)

	verifyList, _ := client.GetNotifications(rpcCtx, &pb.GetNotificationsRequest{UserId: "user_123", Page: 1, Limit: 5})
	for _, n := range verifyList.Notifications {
		if n.Id == targetID {
			if n.IsRead {
				fmt.Printf("TEYIT BASARILI: Bildirim artik 'OKUNDU' olarak gorunuyor.\n")
			} else {
				fmt.Printf("HATA: Bildirim hala 'OKUNMADI' gorunuyor!\n")
			}
		}
	}

	fmt.Println("\nTest Senaryosu Tamamlandi.")
}
