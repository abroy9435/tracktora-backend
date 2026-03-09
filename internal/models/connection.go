package models

import "time"

type Connection struct {
	ID          string    `json:"id"`
	RequesterID string    `json:"requester_id"`
	ReceiverID  string    `json:"receiver_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConnectionRequest struct {
	FriendID string `json:"friend_id"`
}

type UpdateConnectionRequest struct {
	FriendID string `json:"friend_id"`
	Status   string `json:"status"` // 'accepted' or 'rejected'
}
