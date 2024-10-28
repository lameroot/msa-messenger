package user_usercase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
	"github.com/lameroot/msa-messenger/internal/user/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_usecase_AddUserToFriends_mockery(t *testing.T) {
	t.Parallel()

	//Arrange
	var (
		ctx = context.Background()
	)

	type args struct {
		ctx context.Context
		IDUser uuid.UUID
		IDFriend uuid.UUID
	}

	tests := []struct {
		name string
		args args
		want *user_models.AddUserToFriendsResponse
		assertErr assert.ErrorAssertionFunc
		mock func(t *testing.T) user_usecase.PersistentRepository
	}{
		{
			name: "success",
			mock: func(t *testing.T) user_usecase.PersistentRepository {
				repoMock := mocks.NewPersistentRepository(t)
				friendShipStatus := user_models.Pending
				repoMock.EXPECT().
					AddUserToFriends(uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"), uuid.MustParse("4dd068b8-47cf-405b-a051-3a4b38129dbd")).
					Return(friendShipStatus, nil).
					Once()

				return repoMock
			},
			args: args {
				ctx: ctx,
				IDUser: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
				IDFriend: uuid.MustParse("4dd068b8-47cf-405b-a051-3a4b38129dbd"),
			},
			want: &user_models.AddUserToFriendsResponse{
				FriendshipStatus: user_models.Pending,
			},
			assertErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// ARRANGE
			userService := user_usecase.NewUserService(tt.mock(t))

			// ACT
			got, err := userService.AddUserToFriend(tt.args.IDUser, tt.args.IDFriend)

			// ASSERT
			tt.assertErr(t, err)
			if got != nil {
				assert.NotEmpty(t, got.FriendshipStatus)
				assert.Equal(t, user_models.Pending, got.FriendshipStatus)
			}
		})
	}
}