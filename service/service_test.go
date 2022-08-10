package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	mock_storage "github.com/NeverlandMJ/arzon-market/mocks"
	"github.com/NeverlandMJ/arzon-market/pkg/product"
	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/NeverlandMJ/arzon-market/service"
	"github.com/NeverlandMJ/arzon-market/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestService() *service.Service {
	var r storage.Repository
	return service.NewService(r)
}

func NewMockRepo(t *testing.T) *mock_storage.MockRepository {
	ctl := gomock.NewController(t)
	mockRepo := mock_storage.NewMockRepository(ctl)
	return mockRepo
}

func TestService_CreateUser(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		ctx := context.Background()
		user := user.PreSignUpUser{
			Name:        "Hasanova Sunbula",
			Password:    "213",
			PhoneNumber: "1234567",
		}

		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)

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

		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(service.ErrUserExist)

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
			Password:    "123",
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
			Password:    "123",
		}

		mockRepo := NewMockRepo(t)
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
			Password:    "123",
		}

		mockRepo := NewMockRepo(t)
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
			Password:    "123",
		}

		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetUser(ctx, login.PhoneNumber, login.Password).Return(user.User{}, nil)

		s := service.NewService(mockRepo)
		token, err := s.LoginUser(ctx, login)

		require.NotEmpty(t, token)
		assert.Equal(t, nil, err)
	})
}

func TestService_CreateCard(t *testing.T) {
	t.Run("should success", func(t *testing.T) {
		ctx := context.Background()
		testCard := user.PreAddCard{
			CardNumber: "123456789",
			Balance:    1_000_000,
		}

		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddCard(gomock.Any(), gomock.Any()).Return(nil)
		s := service.NewService(mockRepo)
		userCard, err := s.CreateCard(ctx, "ownerID", testCard)
		assert.NotEmpty(t, userCard)
		require.Equal(t, nil, err)
	})
	t.Run("should error", func(t *testing.T) {
		ctx := context.Background()
		testCard := user.PreAddCard{
			CardNumber: "123456789",
			Balance:    1_000_000,
		}

		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddCard(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
		s := service.NewService(mockRepo)
		userCard, err := s.CreateCard(ctx, "ownerID", testCard)
		assert.Empty(t, userCard)
		require.Equal(t, service.ErrServer, err)
	})
}

func TestService_SellProduct(t *testing.T) {
	ctx := context.Background()
	p := product.Product{
		ID:            "dc5444bb-9343-4742-89fc-46f4898d7124",
		Name:          "shaftoli",
		Description:   "shririn shaftoli",
		Quantity:      5,
		Price:         1000,
		OriginalPrice: 500,
		ImageLink:     "link",
		Category:      "meva",
	}
	t.Run("when product doesn't exist", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).Return(product.Product{}, sql.ErrNoRows)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 1, "uID")
		assert.Equal(t, service.ErrProductNotExist, err)
	})
	t.Run("when server error happens inside GetProduct()", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).Return(product.Product{}, fmt.Errorf("server error"))

		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 1, "uID")
		assert.Equal(t, service.ErrServer, err)
	})
	t.Run("when user's card doesn't exist", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(0, sql.ErrNoRows),
		)

		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 1, "uID")
		assert.Equal(t, service.ErrCardNotExist, err)
	})
	t.Run("when server error happens inside GetBalance()", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(0, fmt.Errorf("internal server error")),
		)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 1, "uID")
		assert.Equal(t, service.ErrServer, err)
	})
	t.Run("when quantity exeeded", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(10000, nil),
		)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 6, "uID")
		assert.Equal(t, service.ErrQuantityExceeded, err)
	})
	t.Run("when user's balance not enough", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(100, nil),
		)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 4, "uID")
		assert.Equal(t, service.ErrNotEnoughBalance, err)
	})

	t.Run("when internal error happens inside repo.SellProduct", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(10000, nil),
			mockRepo.EXPECT().SellProduct(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(fmt.Errorf("server error")),
		)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 4, "uID")
		assert.Equal(t, service.ErrServer, err)
	})
	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		gomock.InOrder(
			mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
				Return(p, nil),

			mockRepo.EXPECT().GetBalance(gomock.Any()).
				Return(10000, nil),
			mockRepo.EXPECT().SellProduct(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil),
		)
		s := service.NewService(mockRepo)
		err := s.SellProduct(ctx, "pID", 4, "uID")
		assert.Equal(t, nil, err)
	})
}

func TestService_AllProduct(t *testing.T) {
	p := product.Product{
		ID:            "dc5444bb-9343-4742-89fc-46f4898d7124",
		Name:          "shaftoli",
		Description:   "shririn shaftoli",
		Quantity:      5,
		Price:         1000,
		OriginalPrice: 500,
		ImageLink:     "link",
		Category:      "meva",
	}
	t.Run("when internal server error happens", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().ListProducts(ctx).Return(nil, fmt.Errorf("error"))
		s := service.NewService(mockRepo)
		products, err := s.AllProducts(ctx)

		require.Empty(t, products)
		require.Equal(t, service.ErrServer, err)
	})
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().ListProducts(ctx).Return([]product.Product{p}, nil)
		s := service.NewService(mockRepo)
		products, err := s.AllProducts(ctx)

		require.Equal(t, []product.Product{p}, products)
		require.Equal(t, nil, err)

	})
}

func TestService_GetOneProductInfo(t *testing.T) {
	ctx := context.Background()
	t.Run("when product doesn't exist", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
			Return(product.Product{}, sql.ErrNoRows)

		s := service.NewService(mockRepo)
		p, err := s.GetOneProductInfo(ctx, "id")
		require.Empty(t, p)
		require.Equal(t, service.ErrProductNotExist, err)
	})
	t.Run("when server errorha]ens inside repo.GetProduct", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
			Return(product.Product{}, fmt.Errorf("error"))

		s := service.NewService(mockRepo)
		p, err := s.GetOneProductInfo(ctx, "id")
		require.Empty(t, p)
		require.Equal(t, service.ErrServer, err)
	})
	t.Run("success", func(t *testing.T) {
		want := product.Product{
			ID:            "dc5444bb-9343-4742-89fc-46f4898d7124",
			Name:          "shaftoli",
			Description:   "shririn shaftoli",
			Quantity:      5,
			Price:         1000,
			OriginalPrice: 500,
			ImageLink:     "link",
			Category:      "meva",
		}
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().GetProduct(gomock.Any(), gomock.Any()).
			Return(want, nil)

		s := service.NewService(mockRepo)
		got, err := s.GetOneProductInfo(ctx, "id")
		require.Equal(t, want, got)
		require.Empty(t, err)
	})
}

func TestService_ProductAdd(t *testing.T) {
	ctx := context.Background()
	t.Run("server error", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

		s := service.NewService(mockRepo)
		err := s.ProductAdd(ctx, product.PreAddProduct{})
		require.Equal(t, service.ErrServer, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any()).Return(nil)

		s := service.NewService(mockRepo)
		err := s.ProductAdd(ctx, product.PreAddProduct{})
		require.Equal(t, nil, err)
	})
}

func TestService_ProductsAdd(t *testing.T) {
	ctx := context.Background()
	t.Run("server error", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddProducts(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

		s := service.NewService(mockRepo)
		err := s.ProductsAdd(ctx, []product.PreAddProduct{})
		require.Equal(t, service.ErrServer, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().AddProducts(gomock.Any(), gomock.Any()).Return(nil)

		s := service.NewService(mockRepo)
		err := s.ProductsAdd(ctx, []product.PreAddProduct{})
		require.Equal(t, nil, err)
	})
}

func TestService_UsersList(t *testing.T) {
	ctx := context.Background()
	t.Run("server error", func(t *testing.T) {
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().ListUsers(gomock.Any()).Return(nil, fmt.Errorf("error"))

		s := service.NewService(mockRepo)
		u, err := s.UsersList(ctx)
		require.Equal(t, service.ErrServer, err)
		require.Empty(t, u)
	})

	t.Run("success", func(t *testing.T) {
		want := []user.UserCard{
			{
				FullName:    "Hasanova Sunbula",
				Password:    "hashed_password",
				PhoneNumber: "12345678",
				CardNumber: "8600****",
				Balance: 1000,
			},
		}
		mockRepo := NewMockRepo(t)
		mockRepo.EXPECT().ListUsers(gomock.Any()).Return(want, nil)

		s := service.NewService(mockRepo)
		got, err := s.UsersList(ctx)
		require.Empty(t, err)
		require.Equal(t, want, got)
	})
}
