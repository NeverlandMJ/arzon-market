package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/NeverlandMJ/arzon-market/service"
	"github.com/NeverlandMJ/arzon-market/storage"
	mock_storage "github.com/NeverlandMJ/arzon-market/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestService() *service.Service {
	var r storage.Repository
	return service.NewService(r)
}

func TestService_CreateUser(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		ctx := context.Background()
		user := user.PreSignUpUser{
			Name:        "Hasanova Sunbula",
			Password:    "213",
			PhoneNumber: "1234567",
		}

		ctl := gomock.NewController(t)
		mockRepo := mock_storage.NewMockRepository(ctl)
		gomock.InOrder(
			mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil),
		)

		s := service.NewService(mockRepo)
		err := s.CreateUser(ctx, user)
		assert.Equal(t, nil, err)
	})

	t.Run("with error1", func(t *testing.T) {
		ctx := context.Background()
		user := user.PreSignUpUser{
			Name:        "Hasanova Sunbula",
			Password:    "213",
			PhoneNumber: "1234567",
		}

		ctl := gomock.NewController(t)
		mockRepo := mock_storage.NewMockRepository(ctl)
		gomock.InOrder(
			mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(service.ErrUserExist),
		)

		s := service.NewService(mockRepo)
		err := s.CreateUser(ctx, user)
		assert.Equal(t, service.ErrUserExist, err)
	})
	t.Run("with error2", func(t *testing.T) {
		ctx := context.Background()
		user := user.PreSignUpUser{
			Name:        "",
			Password:    "",
			PhoneNumber: "",
		}

		s := NewTestService()
		err := s.CreateUser(ctx, user)
		assert.Equal(t, service.ErrInvalidUser, err)
	})

}

func TestService_LoginUser(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		ctx := context.Background()
		login := user.PreLoginUser{
			PhoneNumber: "",
			Password: "123",
		}
		s := NewTestService()
		token, err := s.LoginUser(ctx, login)

		require.Empty(t, token)
		assert.Equal(t, service.ErrInvalidUser, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		ctx := context.Background()
		login := user.PreLoginUser{
			PhoneNumber: "1234567",
			Password: "123",
		}

		ctl := gomock.NewController(t)
		mockRepo := mock_storage.NewMockRepository(ctl)
		mockRepo.EXPECT().GetUser(ctx, login.PhoneNumber, login.Password).Return(user.User{}, sql.ErrNoRows)

		s := service.NewService(mockRepo)
		token, err := s.LoginUser(ctx, login)

		require.Empty(t, token)
		assert.Equal(t, service.ErrUserNotExist, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		ctx := context.Background()
		login := user.PreLoginUser{
			PhoneNumber: "1234567",
			Password: "123",
		}

		ctl := gomock.NewController(t)
		mockRepo := mock_storage.NewMockRepository(ctl)
		mockRepo.EXPECT().GetUser(ctx, login.PhoneNumber, login.Password).Return(user.User{}, fmt.Errorf("server error"))

		s := service.NewService(mockRepo)
		token, err := s.LoginUser(ctx, login)

		require.Empty(t, token)
		assert.Equal(t, service.ErrServer, err)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		login := user.PreLoginUser{
			PhoneNumber: "1234567",
			Password: "123",
		}

		ctl := gomock.NewController(t)
		mockRepo := mock_storage.NewMockRepository(ctl)
		mockRepo.EXPECT().GetUser(ctx, login.PhoneNumber, login.Password).Return(user.User{}, nil)

		s := service.NewService(mockRepo)
		token, err := s.LoginUser(ctx, login)

		require.NotEmpty(t, token)
		assert.Equal(t, nil, err)
	})
}

func TestService_CreateCard(t *testing.T) {
	
}