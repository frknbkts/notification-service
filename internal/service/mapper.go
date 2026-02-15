package service

import (
	"github.com/frknbkts/notification-service/internal/core/domain"
	"github.com/frknbkts/notification-service/pkg/pb"
)

func ProtoToDomain(req *pb.SendNotificationRequest) *domain.Notification {
	return &domain.Notification{
		UserID:      req.UserId,
		SenderID:    req.SenderId,
		Type:        req.Type.String(),
		Title:       req.Title,
		Body:        req.Body,
		ReferenceID: req.ReferenceId,
		IsRead:      false,
	}
}

func DomainToProto(n *domain.Notification) *pb.Notification {
	var notifType pb.NotificationType
	if val, ok := pb.NotificationType_value[n.Type]; ok {
		notifType = pb.NotificationType(val)
	} else {
		notifType = pb.NotificationType_UNKNOWN
	}

	return &pb.Notification{
		Id:          n.ID,
		UserId:      n.UserID,
		SenderId:    n.SenderID,
		Type:        notifType,
		Title:       n.Title,
		Body:        n.Body,
		ReferenceId: n.ReferenceID,
		IsRead:      n.IsRead,
		CreatedAt:   n.CreatedAt,
	}
}
