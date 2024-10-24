package user_usecase

import (
	"errors"
	"log"

	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
)

type UserService struct {
	persistentRepository PersistentRepository
}

func NewUserService(persistentRepository PersistentRepository) *UserService {
	return &UserService{
		persistentRepository: persistentRepository,
	}
}

func (u *UserService) AddUserToFriend(IDUser uuid.UUID, IDFriend uuid.UUID) (*user_models.AddUserToFriendsResponse, error) {
	friendshipStatus, error := u.persistentRepository.AddUserToFriends(IDUser, IDFriend)
	if error != nil {
		log.Default().Printf("Error add user " + IDFriend.String() + " to user " + IDUser.String() + ": " + error.Error())
		return nil, error
	}
	addUserToFriendsResponse := &user_models.AddUserToFriendsResponse{
		FriendshipStatus: friendshipStatus,
	}
	return addUserToFriendsResponse, error
}

func (u *UserService) DeleteUserFromFriends(IDUser uuid.UUID, IDFriend uuid.UUID) error {
	error := u.persistentRepository.DeleteUserFromFriends(IDUser, IDFriend)
	if error != nil {
		log.Default().Printf("Error delete user " + IDFriend.String() + " from user " + IDUser.String() + " : " + error.Error())
		return errors.New("Error delete user " + IDFriend.String() + " from user " + IDUser.String())
	}
	return nil
}

func (u *UserService) AcceptFriend(IDUser uuid.UUID, IDFriend uuid.UUID) error {
	error := u.persistentRepository.AcceptFriend(IDUser, IDFriend)
	if error != nil {
		log.Default().Printf("Error accept friendship " + IDFriend.String() + " from user " + IDUser.String() + ": " + error.Error())
		return errors.New("Error accept friendship " + IDFriend.String() + " from user " + IDUser.String())
	}
	return nil
}

func (u *UserService) RejectFriend(IDUser uuid.UUID, IDFriend uuid.UUID) error {
	error := u.persistentRepository.RejectFriend(IDUser, IDFriend)
	if error != nil {
		log.Default().Printf("Error reject friendship " + IDFriend.String() + " from user " + IDUser.String() + ": " + error.Error())
		return errors.New("Error reject friendship " + IDFriend.String() + " from user " + IDUser.String())
	}
	return nil
}

func (u *UserService) GetFriends(IDUser uuid.UUID) ([]user_models.Friend, error) {
	friends, error := u.persistentRepository.GetFriends(IDUser)
	if error != nil {
		log.Default().Printf("Error get friends from user " + IDUser.String() + ": " + error.Error())
		return []user_models.Friend{}, errors.New("Error get friends from user " + IDUser.String())
	}
	return friends, nil
}
