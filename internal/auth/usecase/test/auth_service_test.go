package auth_usecase_test

import (
	"log"
	"testing"

	"github.com/google/uuid"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
	"github.com/lameroot/msa-messenger/internal/auth/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_Resgister_minimock(t *testing.T) {
	t.Parallel()

	type args struct {
		email string
		password string
		nickname string
	}

	tests := []struct {
		name string
		args args
		want auth_models.User
		wantErr assert.ErrorAssertionFunc
		mock func(t *testing.T) auth_usecase.PersistentRepository
	}{
		{
			name : "Success",
			args: args{
				email: "test@mail.ru",
				password: "test",
				nickname: "test",
			},
			want: auth_models.User{
				ID: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
				Email: "test@mail.ru",
				Password: "test",
				Nickname: "test",
				AvatarURL: "test.ru",
				Info: "info",
			},
			wantErr: assert.NoError,
			mock: func(t *testing.T) auth_usecase.PersistentRepository {
				repoMock := mocks.NewPersistentRepositoryMock(t)
				
				repoMock.GetUserByEmailMock.
					Expect("test@mail.ru").
					Times(1).
					Return(nil, nil)
					
				repoMock.GetUserByNicknameMock.
					Expect("test").
					Times(1).
					Return(nil, nil)
				
				repoMock.SaveUserMock.
					Times(1).
					Return(&auth_models.User{
						ID: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
						Email: "test@mail.ru",
						Password: "test",
						Nickname: "test",
						AvatarURL: "test.ru",
						Info: "info",
					}, nil)
					
				return repoMock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authService := auth_usecase.NewAuthUseCase(mocks.NewTokenRepositoryMock(t), tt.mock(t))
			got, err := authService.Register(tt.args.email, tt.args.password, tt.args.nickname)
			tt.wantErr(t, err)

			assert.Equal(t, uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"), got.ID)
		})
	}
	
}

func Test_Authenticate_minimock(t *testing.T) {
	t.Parallel()

	type args struct {
		email string
		password string
	}
	type Deps struct {
		tokenRepository auth_usecase.TokenRepository
		persistentRepository auth_usecase.PersistentRepository
	}

	tests := []struct{
		name string
		args args
		want auth_models.AuthResponse
		wantErr assert.ErrorAssertionFunc
		mock func(t *testing.T) Deps
	}{
		{
			name: "success",
			args: args{
				email: "test@mail.ru",
				password: "testtest",
			},
			want: auth_models.AuthResponse{
				User: auth_models.User{
					ID: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
				},
				Token: auth_models.Token{
					AccessToken: "access_token",
					RefreshToken: "refresh_token",
				},
			},
			wantErr: assert.NoError,
			mock: func(t *testing.T) Deps {
				tokenRepo := mocks.NewTokenRepositoryMock(t)
				persRepo := mocks.NewPersistentRepositoryMock(t)

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testtest"), bcrypt.DefaultCost)
				persRepo.GetUserByEmailMock.
					Expect("test@mail.ru").
					Times(1).
					Return(&auth_models.User{
						ID: uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"),
						Password: string(hashedPassword),
					}, nil)

				tokenRepo.CreateAccessTokenMock.
					Times(1).
					Return("access_token", nil)

				tokenRepo.CreateRefreshTokenMock.
					Return("refresh token", nil)
					
				return Deps{
					tokenRepository: tokenRepo,
					persistentRepository: persRepo,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := tt.mock(t)

			authService := auth_usecase.NewAuthUseCase(deps.tokenRepository, deps.persistentRepository)
			log.Default().Println(authService)
			got, err := authService.Authenticate(tt.args.email, tt.args.password)			
			log.Default().Println(got)
			log.Default().Println(err)
			// tt.wantErr(t, err)

			// assert.Equal(t, "access_token", got.Token.AccessToken)
			// assert.Equal(t, "refresh_token", got.Token.RefreshToken)
			// assert.Equal(t, uuid.MustParse("99617edb-249c-4c04-9889-94cdf6eb4e00"), got.User.ID)
		})
	}
}