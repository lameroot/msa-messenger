package messaging_notification

import (
	"log"

	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
)

type InMemmoryNotificationService struct {
	messages []messaging_models.Message
}

func NewInMemmoryNotificationService() messaging_usecase.NotificationService {
	return &InMemmoryNotificationService{}
}

func (s *InMemmoryNotificationService) SendNotification(IDFrom uuid.UUID, IDTo uuid.UUID, SentTime int64, ShortMessage string) (string, error) {
	message := messaging_models.Message{}
	message.IDSender = IDFrom
	message.IDReceiver = IDTo
	message.Message = ShortMessage
	message.SentTime = SentTime

	s.messages = append(s.messages, message)

	return "OK", nil
}

func (s *InMemmoryNotificationService) Close() {
	log.Default().Println("Close InMemmoryNotificationService")
}
