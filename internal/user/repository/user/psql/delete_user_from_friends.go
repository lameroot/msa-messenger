package user_repository_psql

import "github.com/google/uuid"

func (s *PosgresUserRepository) DeleteUserFromFriends(IDUser uuid.UUID, IDFriend uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM friendships WHERE user_id = $1 AND friend_id = $2",
		IDUser, IDFriend)
	if err != nil {
		return err
	}
	return nil
}
