package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/frknbkts/notification-service/internal/core/domain"
	"github.com/frknbkts/notification-service/internal/core/ports"
	"github.com/google/uuid"
)

type CouchbaseRepository struct {
	cluster    *gocb.Cluster
	collection *gocb.Collection
	scope      *gocb.Scope
}

func NewCouchbaseRepository(connectionString, username, password, bucketName string) (ports.NotificationRepository, error) {
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("couchbase baglanti hatasi: %v", err)
	}

	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(20*time.Second, nil)
	if err != nil {
		return nil, fmt.Errorf("couchbase bucket hazir degil: %v", err)
	}

	scope := bucket.Scope("app")
	collection := scope.Collection("notifications")

	return &CouchbaseRepository{
		cluster:    cluster,
		collection: collection,
		scope:      scope,
	}, nil
}

func (r *CouchbaseRepository) Save(ctx context.Context, n *domain.Notification) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	docKey := fmt.Sprintf("notification::%s", n.ID)

	_, err := r.collection.Insert(docKey, n, &gocb.InsertOptions{
		Context: ctx,
		Expiry:  30 * 24 * 60 * 60 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("kayit basarisiz: %v", err)
	}
	return nil
}

func (r *CouchbaseRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Notification, int, error) {
	query := `
		SELECT t.* FROM notifications as t 
		WHERE t.user_id = $1 
		ORDER BY t.created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.scope.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{userID, limit, offset},
		Context:              ctx,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("sorgu hatasi: %v", err)
	}

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Row(&n); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, &n)
	}

	return notifications, len(notifications), nil
}

func (r *CouchbaseRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	docKey := fmt.Sprintf("notification::%s", notificationID)

	ops := []gocb.MutateInSpec{
		gocb.UpsertSpec("is_read", true, nil),
	}

	_, err := r.collection.MutateIn(docKey, ops, &gocb.MutateInOptions{
		Context: ctx,
	})

	return err
}

func (r *CouchbaseRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*) as count 
		FROM notifications 
		WHERE user_id = $1 AND is_read = false`

	rows, err := r.scope.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{userID},
		Context:              ctx,
	})
	if err != nil {
		return 0, err
	}

	var result struct {
		Count int `json:"count"`
	}

	if rows.Next() {
		rows.Row(&result)
	}
	return result.Count, nil
}
