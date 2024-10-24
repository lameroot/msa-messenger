package messaging_usecase

import (
	"time"

	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
)

type MessagingService struct {
	persistentRepository PersistentRepository
	notificationService  NotificationService
}

func NewMessagingService(persistentRepository PersistentRepository, notificationService NotificationService) *MessagingService {
	return &MessagingService{
		persistentRepository: persistentRepository,
		notificationService:  notificationService,
	}
}

func (s *MessagingService) SendMessage(IDUser uuid.UUID, sendMessageRequest messaging_models.SendMessageRequest) (messaging_models.SendMessageResponse, error) {
	err := s.persistentRepository.SaveMessage(IDUser, sendMessageRequest)
	if err != nil {
		return messaging_models.SendMessageResponse{}, err
	}
	_, err = s.notificationService.SendNotification(IDUser, sendMessageRequest.IDFriend, sendMessageRequest.SentTime, "You have one new message from "+IDUser.String())
	if err != nil {
		return messaging_models.SendMessageResponse{}, err
	}
	return messaging_models.SendMessageResponse{Success: true, DeliveryTime: time.Now().Unix()}, nil
}

func (s *MessagingService) GetMessages(IDUser uuid.UUID, getMessagesRequest messaging_models.GetMessagesRequest) ([]messaging_models.Message, error) {
	messages, err := s.persistentRepository.GetMessages(IDUser, getMessagesRequest.CountMessages)
	if err != nil {
		return []messaging_models.Message{}, err
	}
	return messages, nil
}
