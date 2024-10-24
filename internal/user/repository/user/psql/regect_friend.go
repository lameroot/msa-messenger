package user_repository_psql

import "github.com/google/uuid"

func (s *PosgresUserRepository) RejectFriend(IDUser uuid.UUID, IDFriend uuid.UUID) error {
	_, err := s.db.Exec("UPDATE friendships SET status = 'rejected' WHERE user_id = $1 and friend_id = $2",
		IDUser, IDFriend)
	if err != nil {
		return err
	}
	return nil
}
