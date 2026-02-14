package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/frknbkts/notification-service/internal/core/ports"
	"github.com/frknbkts/notification-service/pkg/pb"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer 
	repo                                      ports.NotificationRepository
}

func NewNotificationService(repo ports.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	if req.UserId == "" || req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id ve title bos olamaz")
	}

	notification := ProtoToDomain(req)
	
	notification.ID = uuid.New().String()
	notification.CreatedAt = time.Now().Unix()

	err := s.repo.Save(ctx, notification)
	if err != nil {
		return &pb.SendNotificationResponse{Success: false}, status.Errorf(codes.Internal, "kayit hatasi: %v", err)
	}

	return &pb.SendNotificationResponse{
		Id:      notification.ID,
		Success: true,
	}, nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10 
	}
	
	offset := (req.Page - 1) * req.Limit
	if offset < 0 {
		offset = 0
	}

	notifications, totalCount, err := s.repo.GetByUserID(ctx, req.UserId, int(req.Limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "veri okuma hatasi: %v", err)
	}

	var pbNotifications []*pb.Notification
	for _, n := range notifications {
		pbNotifications = append(pbNotifications, DomainToProto(n))
	}

	return &pb.GetNotificationsResponse{
		Notifications: pbNotifications,
		TotalCount:    int32(totalCount),
		CurrentPage:   req.Page,
	}, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	err := s.repo.MarkAsRead(ctx, req.NotificationId, req.UserId)
	if err != nil {
		return &pb.MarkAsReadResponse{Success: false}, status.Errorf(codes.Internal, "guncelleme hatasi: %v", err)
	}
	return &pb.MarkAsReadResponse{Success: true}, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := s.repo.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sayim hatasi: %v", err)
	}
	return &pb.GetUnreadCountResponse{Count: int32(count)}, nil
}