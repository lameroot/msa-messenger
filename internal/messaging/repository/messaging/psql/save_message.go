package messaging_repository_psql

import (
	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
)

func (p *PostgresMessagingRepository) SaveMessage(IDUser uuid.UUID, messageRequest messaging_models.SendMessageRequest) error {
	var ID string = uuid.New().String()
	_, err := p.db.Exec("INSERT INTO messages (id, sender_id, receiver_id, content) values ($1, $2, $3, $4)",
		ID, IDUser.String(), messageRequest.IDFriend.String(), messageRequest.Message)
	if err != nil {
		return err
	}
	return nil
}
