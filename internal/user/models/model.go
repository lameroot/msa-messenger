package user_models

import "github.com/google/uuid"

// AddUserToFriendsRequest is the request model for adding a user to friend
type AddUserToFriendsRequest struct {
	FriendUserID uuid.UUID `json:"friend_id"`
}

// AddUserToFriendsResponse is the response model for adding a user to friend
type AddUserToFriendsResponse struct {
	FriendshipStatus FriendshipStatus `json:"friendship_status"`
}

// DeleteUserFromFriendsRequest is the request model for deleting a user from friend
type DeleteUserFromFriendsRequest struct {
	FriendUserID uuid.UUID `json:"friend_id"`
}

// AddUserToFriendRequest is the request model for adding a user to friend
type AcceptFriendRequest struct {
	FriendUserID uuid.UUID `json:"friend_id"`
}

// DeleteUserFromFriendRequest is the request model for deleting a user from friend
type RejectFriendRequest struct {
	FriendUserID uuid.UUID `json:"friend_id"`
}

// GetFriendsRequest is the request model for getting friends
type GetFriendsRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

// Friend is the model for a friend
type Friend struct {
	ID               uuid.UUID        `json:"friend_id"`
	FriendshipStatus FriendshipStatus `json:"friendship_status"`
}

// FriendshipStatus is the model for friendship status
type FriendshipStatus string

const (
	Pending  FriendshipStatus = "pending"
	Accepted FriendshipStatus = "accepted"
	Rejected FriendshipStatus = "rejected"
)

// Pending is the model for pending friendship request
type GetFriendsResponse struct {
	Friends []Friend `json:"friends"`
}
