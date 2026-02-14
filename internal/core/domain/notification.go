package domain

type Notification struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	SenderID    string `json:"sender_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	ReferenceID string `json:"reference_id"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   int64  `json:"created_at"`
}