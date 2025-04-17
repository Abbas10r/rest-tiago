package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, followingUserID int64) error {
	query := `
		INSERT INTO followers (following_id, follower_id) VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followingUserID, followerID)
	if err != nil {
		return err
	}

	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, followingUserID int64) error {
	query := `
		DELETE FROM followers 
		WHERE following_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followingUserID, followerID)
	if err != nil {
		return err
	}

	return nil
}
