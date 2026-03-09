package repository

import (
	"context"
	"errors"
	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
)

// SendInviteByID: User A (requesterID) invites User B (receiverID)
func SendInviteByID(requesterID, receiverID string) error {
	if requesterID == receiverID {
		return errors.New("you cannot add yourself as a friend")
	}

	// 1. Insert: requester is 'user_id', receiver is 'friend_id'
	query := `INSERT INTO friends (user_id, friend_id, status) VALUES ($1, $2, 'pending')`
	_, err := database.DB.Exec(context.Background(), query, requesterID, receiverID)
	if err != nil {
		return errors.New("request already exists or users are already connected")
	}
	return nil
}

// GetPendingRequests: Show User B all the User As who invited them
func GetPendingRequests(userID string) ([]map[string]interface{}, error) {
	// Changed u.first_name to u.username to match your users table
	query := `
		SELECT 
			u.id AS sender_id, 
			u.username AS sender_name, 
			u.email AS sender_email
		FROM friends f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = $1 AND f.status = 'pending'`

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var sID, sName, sEmail string
		if err := rows.Scan(&sID, &sName, &sEmail); err != nil {
			return nil, err
		}
		requests = append(requests, map[string]interface{}{
			"sender_id":    sID,
			"sender_name":  sName,
			"sender_email": sEmail,
		})
	}
	return requests, nil
}

// UpdateFriendStatus: User B accepts a request from User A (senderID)
func UpdateFriendStatus(userID, senderID, newStatus string) error {
	query := `UPDATE friends SET status = $1 WHERE friend_id = $2 AND user_id = $3`
	result, err := database.DB.Exec(context.Background(), query, newStatus, userID, senderID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("no pending request found from this user")
	}
	return nil
}

// GetFriendStats: Verify 'accepted' status before showing data
func GetFriendStats(userID, friendID string) (*models.ApplicationStats, error) {
	var isSharing bool
	var isAccepted bool

	query := `
		SELECT 
			u.share_stats,
			EXISTS (
				SELECT 1 FROM friends 
				WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
				AND status = 'accepted'
			) as is_accepted
		FROM users u
		WHERE u.id = $2`

	err := database.DB.QueryRow(context.Background(), query, userID, friendID).Scan(&isSharing, &isAccepted)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !isAccepted {
		return nil, errors.New("access denied: you must be accepted friends to view stats")
	}

	// Now it only blocks if a user explicitly went into settings and turned it OFF
	if !isSharing {
		return nil, errors.New("this user has chosen to keep their stats private")
	}

	return GetApplicationStats(friendID)
}

// GetAcceptedFriends returns a list of all users you are officially connected with
func GetAcceptedFriends(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			u.id, u.username, u.email, u.share_stats
		FROM users u
		JOIN friends f ON (f.user_id = u.id OR f.friend_id = u.id)
		WHERE (f.user_id = $1 OR f.friend_id = $1) 
		  AND f.status = 'accepted'
		  AND u.id != $1` // Don't include yourself in the list

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []map[string]interface{}
	for rows.Next() {
		var id, name, email string
		var sharing bool
		if err := rows.Scan(&id, &name, &email, &sharing); err != nil {
			return nil, err
		}
		friends = append(friends, map[string]interface{}{
			"friend_id":   id,
			"username":    name,
			"email":       email,
			"share_stats": sharing,
		})
	}
	return friends, nil
}
