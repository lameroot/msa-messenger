package messaging_models

import "github.com/google/uuid"

type SendMessageRequest struct {
	IDFriend uuid.UUID `json:"friend_id"`
	Message  string    `json:"message"`
	SentTime int64     `json:"sent_time"`
}

type SendMessageResponse struct {
	Success      bool  `json:"success"`
	DeliveryTime int64 `json:"delivery_time"`
}

type GetMessagesRequest struct {
	IDFriend      uuid.UUID `json:"friend_id"`
	CountMessages int32     `json:"count_messages"`
}

type PullMessagesResponse struct {
	Success bool   `json:"success"`
	Message []byte `json:"messages"`
}

type Message struct {
	IDSender   uuid.UUID `json:"sender_id"`
	IDReceiver uuid.UUID `json:"receiver_id"`
	Message    string    `json:"message"`
	Read       bool      `json:"read"`
	SentTime   int64     `json:"sent_time"`
}
