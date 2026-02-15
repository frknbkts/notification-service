package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/frknbkts/notification-service/internal/core/domain"
	"github.com/frknbkts/notification-service/pkg/pb"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, n *domain.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Notification, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*domain.Notification), args.Int(1), args.Error(2)
}

func (m *MockRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	args := m.Called(ctx, notificationID, userID)
	return args.Error(0)
}

func (m *MockRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func TestSendNotification_Success(t *testing.T) {

	mockRepo := new(MockRepository)
	service := NewNotificationService(mockRepo)

	ctx := context.Background()
	req := &pb.SendNotificationRequest{
		UserId:   "user_123",
		SenderId: "user_999",
		Type:     pb.NotificationType_LIKE,
		Title:    "Test Bildirimi",
		Body:     "Merhaba",
	}

	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	resp, err := service.SendNotification(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Id)

	mockRepo.AssertExpectations(t)
}

func TestGetNotifications_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewNotificationService(mockRepo)
	ctx := context.Background()

	fakeNotifications := []*domain.Notification{
		{ID: "1", Title: "Notif 1", IsRead: false},
		{ID: "2", Title: "Notif 2", IsRead: true},
	}

	mockRepo.On("GetByUserID", ctx, "user_123", 10, 0).Return(fakeNotifications, 2, nil)

	req := &pb.GetNotificationsRequest{UserId: "user_123", Page: 1, Limit: 10}
	resp, err := service.GetNotifications(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Notifications))
	assert.Equal(t, "Notif 1", resp.Notifications[0].Title)

	mockRepo.AssertExpectations(t)
}
func TestMarkAsRead_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewNotificationService(mockRepo)
	ctx := context.Background()
	notifID := "bildirim_123"
	userID := "user_123"

	mockRepo.On("MarkAsRead", ctx, notifID, userID).Return(nil)

	req := &pb.MarkAsReadRequest{NotificationId: notifID, UserId: userID}
	resp, err := service.MarkAsRead(ctx, req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	mockRepo.AssertExpectations(t)
}

func TestGetUnreadCount_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewNotificationService(mockRepo)
	ctx := context.Background()
	userID := "user_123"

	mockRepo.On("GetUnreadCount", ctx, userID).Return(5, nil)

	req := &pb.GetUnreadCountRequest{UserId: userID}
	resp, err := service.GetUnreadCount(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int32(5), resp.Count)
	mockRepo.AssertExpectations(t)
}
