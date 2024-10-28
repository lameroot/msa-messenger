package user_repository_psql

import (
	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
)

func (s *PosgresUserRepository) GetFriends(IDUser uuid.UUID) ([]user_models.Friend, error) {
	rows, err := s.db.Query("SELECT friend_id, status FROM friendships WHERE user_id = $1", IDUser.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []user_models.Friend
	for rows.Next() {
		var friend user_models.Friend
		var status string
		err := rows.Scan(&friend.ID, &status)
		if err != nil {
			return nil, err
		}
		friend.FriendshipStatus = user_models.FriendshipStatus(status)
		friends = append(friends, friend)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}
