package messaging_repository_psql

import (
	"log"
	"time"

	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
)

func (p *PostgresMessagingRepository) GetMessages(IDUser uuid.UUID, LastCount int32) ([]messaging_models.Message, error) {
	rows, err := p.db.Query("SELECT id, sender_id, receiver_id, content, is_read, created_at FROM messages where sender_id = $1 ORDER BY created_at DESC LIMIT $2 ",
		IDUser.String(), LastCount)
	if err != nil {
		return []messaging_models.Message{}, err
	}
	defer rows.Close()
	var messages []messaging_models.Message
	for rows.Next() {
		var message messaging_models.Message
		var id, senderId, receiverId string
		var content string
		var isRead bool
		var createdAt time.Time
		err = rows.Scan(&id, &senderId, &receiverId, &content, &isRead, &createdAt)
		if err != nil {
			log.Default().Printf("Error get messages: %s", err)
			continue
		}
		message.IDSender = IDUser
		message.IDReceiver, _ = uuid.Parse(receiverId)
		message.Message = content
		message.Read = isRead
		message.SentTime = createdAt.Unix()

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return []messaging_models.Message{}, err
	}
	return messages, nil
}
