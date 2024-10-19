package user_repository_psql

import (
	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
)

func (s *PosgresUserRepository) AddUserToFriends(IDUser uuid.UUID, IDFriend uuid.UUID) (user_models.FriendshipStatus, error) {
	_, err := s.db.Exec("INSERT INTO friendships (id, user_id, friend_id, status) VALUES ($1, $2, $3, $4)",
		uuid.New().String(), IDUser, IDFriend, user_models.Pending)
	if err != nil {
		return user_models.Rejected, err
	}
	return user_models.Pending, nil
}
