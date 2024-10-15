package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Friendship struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FriendID  int64     `json:"friend_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func AddFriend(db *sql.DB, userID, friendID int64) error {
	_, err := db.Exec(`
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (user_id, friend_id) DO NOTHING
	`, userID, friendID)
	return err
}

func RemoveFriend(db *sql.DB, userID, friendID int64) error {
	_, err := db.Exec(`
		DELETE FROM friendships
		WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, userID, friendID)
	return err
}

func GetFriends(db *sql.DB, userID int64) ([]Friendship, error) {
	rows, err := db.Query(`
		SELECT id, user_id, friend_id, status, created_at, updated_at
		FROM friendships
		WHERE user_id = $1 OR friend_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []Friendship
	for rows.Next() {
		var f Friendship
		err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	return friendships, nil
}

func AcceptFriendRequest(db *sql.DB, userID, friendID int64) error {
	_, err := db.Exec(`
		UPDATE friendships
		SET status = 'accepted', updated_at = NOW()
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`, userID, friendID)
	return err
}

func RejectFriendRequest(db *sql.DB, userID, friendID int64) error {
	_, err := db.Exec(`
		DELETE FROM friendships
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`, userID, friendID)
	return err
}
