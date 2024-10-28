package messaging_usecase_test

import (
	"testing"

	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
	"github.com/lameroot/msa-messenger/internal/messaging/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_usecase_SendMessage_mockery(t *testing.T) {
	t.Parallel()

	type args struct {
		IDUser uuid.UUID
		SendMessageRequest *messaging_models.SendMessageRequest
	}

	type Deps struct {
		persistentRepository messaging_usecase.PersistentRepository
		notificationService messaging_usecase.NotificationService
	}

	tests := []struct{
		name string
		args args
		want *messaging_models.SendMessageResponse
		assertErr assert.ErrorAssertionFunc
		mock func(t *testing.T) Deps
	} {
		{
			name: "success",
			args: args {
				IDUser: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
				SendMessageRequest: &messaging_models.SendMessageRequest{
					IDFriend: uuid.MustParse("4dd068b8-47cf-405b-a051-3a4b38129dbd"),
					Message: "test message",
					SentTime: int64(1729802370),
				},
			},
			assertErr: assert.NoError,
			want: &messaging_models.SendMessageResponse{
				Success: true,
				DeliveryTime: int64(1729802370),
			},
			mock: func (t *testing.T) Deps {
				repoMock := mocks.NewPersistentRepository(t)
				notifyMock := mocks.NewNotificationService(t)
	
				IDUser := uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00")
				IDFriend := uuid.MustParse("4dd068b8-47cf-405b-a051-3a4b38129dbd")
				sendMessageRequest := messaging_models.SendMessageRequest{
					IDFriend: IDFriend,
					Message: "test message",
					SentTime: int64(1729802370),
				}
				repoMock.EXPECT().
					SaveMessage(IDUser, sendMessageRequest).
					Return(nil).
					Once()
				
				notifyMock.EXPECT().
					SendNotification(IDUser, IDFriend, int64(1729802370), "You have one new message from " + IDUser.String()).
					Return("ok", nil).
					Once()
	
				return Deps{
					persistentRepository: repoMock,
					notificationService: notifyMock,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := tt.mock(t)
			// ARRANGE
			messaging_service := messaging_usecase.NewMessagingService(deps.persistentRepository, deps.notificationService)

			// ACT
			got, err := messaging_service.SendMessage(tt.args.IDUser, *tt.args.SendMessageRequest)

			// ASSERT
			tt.assertErr(t, err)
			assert.True(t, got.Success)
			assert.NotEmpty(t, got.DeliveryTime)
		})
	}
	
}