package user_usecase

import (
	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
)

type (
	PersistentRepository interface {
		AddUserToFriends(IDUser uuid.UUID, IDFriend uuid.UUID) (user_models.FriendshipStatus, error)
		DeleteUserFromFriends(IDUser uuid.UUID, IDFriend uuid.UUID) error
		AcceptFriend(IDUser uuid.UUID, IDFriend uuid.UUID) error
		RejectFriend(IDUser uuid.UUID, IDFriend uuid.UUID) error
		GetFriends(IDUser uuid.UUID) ([]user_models.Friend, error)
	}
)
