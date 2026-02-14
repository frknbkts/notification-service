package ports

import (
	"context"
	"github.com/frknbkts/notification-service/internal/core/domain"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Notification, int, error)
	MarkAsRead(ctx context.Context, notificationID string, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}