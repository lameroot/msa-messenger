package messaging_usecase

import (
	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
)

type (
	PersistentRepository interface {
		SaveMessage(IDUser uuid.UUID, messageRequest messaging_models.SendMessageRequest) error
		GetMessages(IDUser uuid.UUID, LastCount int32) ([]messaging_models.Message, error)
		Close()
	}
	NotificationService interface {
		SendNotification(IDFrom uuid.UUID, IDTo uuid.UUID, SentTime int64, ShortMessage string) (string, error)
		Close()
	}
)