package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/NeverlandMJ/arzon-market/pkg/product"
	"github.com/NeverlandMJ/arzon-market/pkg/store"
	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/NeverlandMJ/arzon-market/storage"
	"github.com/dgrijalva/jwt-go"
)

var ErrUserExist = errors.New("user already exist")
var ErrInvalidUser = errors.New("empty faild")
var ErrUserNotExist = errors.New("user doesn't exist")
var ErrServer = errors.New("server error")
var ErrProductNotExist = errors.New("product doesn't exist")
var ErrQuantityExceeded = errors.New("quantity exceeded")
var ErrCardNotExist = errors.New("card isn't set")
var ErrNotEnoughBalance = errors.New("balance isn't enough")

type Handler interface {
	CreateUser(ctx context.Context, tempUser user.PreSignUpUser)  error
	LoginUser(ctx context.Context, tempUser user.PreLoginUser) (string, error)
	CreateCard(ctx context.Context, ownerID string, card user.PreAddCard) (user.Card, error)
	SellProduct(ctx context.Context, productName string, quantity int, uID string) error
	AllProducts(ctx context.Context) ([]product.Product, error)
	GetOneProductInfo(ctx context.Context, name string) (product.Product, error)
	ProductAdd(ctx context.Context, p product.PreAddProduct) error
	ProductsAdd(ctx context.Context, ps []product.PreAddProduct) error
	UsersList(ctx context.Context) ([]user.UserCard, error)
}

type Service struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, tempUser user.PreSignUpUser)  error {

	newUser, err := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.PhoneNumber,
	)
	if err != nil {
		return ErrInvalidUser
	}

	err = s.repo.AddUser(ctx, *newUser)

	if err != nil {
		// log.Println("Service(): ", err)
		return ErrUserExist
	}

	return  nil
}

func (s *Service) LoginUser(ctx context.Context, tempUser user.PreLoginUser) (string, error) {
	if tempUser.PhoneNumber == "" || tempUser.Password == "" {
		return "", ErrInvalidUser
	}

	u, err := s.repo.GetUser(ctx, tempUser.PhoneNumber, tempUser.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotExist
		} else {
			log.Println("Service(): ", err)
			return "", ErrServer
		}
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:          u.ID,
		PhoneNumber: u.PhoneNumber,
		IsAdmin:     u.IsAdmin,
		FullName:    u.FullName,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", ErrServer
	}

	return tokenString, nil
}

func (s *Service) CreateCard(ctx context.Context, ownerID string, card user.PreAddCard) (user.Card, error) {
	newCard := user.NewCard(card.CardNumber, card.Balance, ownerID)

	err := s.repo.AddCard(ctx, *newCard)
	if err != nil {
		log.Println("Service(): ", err)
		return user.Card{}, ErrServer
	}

	return *newCard, nil
} 

func (s *Service) SellProduct(ctx context.Context, proID string, quantity int, uID string) error {
	gotProduct, err := s.repo.GetProduct(ctx, proID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotExist
		} else {
			log.Println(err)
			return ErrServer
		}
	}

	userBalance, err := s.repo.GetBalance(uID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCardNotExist
		} else {
			log.Println(err)
			return ErrServer
		}
	}

	sales, soldProduct, err := store.Sell(gotProduct, quantity, userBalance, uID)

	if err != nil {
		if errors.Is(err, store.ErrNotEnoughQuantity) {
			return ErrQuantityExceeded
		}else {
			return ErrNotEnoughBalance
		}
	}

	err = s.repo.SellProduct(ctx, sales, soldProduct)
	if err != nil {
		log.Println(err)
		return ErrServer
	}

	return nil

}

func (s *Service) AllProducts(ctx context.Context) ([]product.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		log.Println(err)
		return nil, ErrServer
	}

	return products, nil
}

func (s *Service) GetOneProductInfo(ctx context.Context, id string) (product.Product, error) {
	p, err := s.repo.GetProduct(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return product.Product{}, ErrProductNotExist
		} else {
			log.Println(err)
			return product.Product{}, ErrServer
		}
	}

	return p, nil
}

func (s *Service) ProductAdd(ctx context.Context, p product.PreAddProduct) error {
	product := product.New(p)
	err := s.repo.AddProduct(ctx, *product)
	if err != nil {
		log.Println(err)
		return ErrServer
	}
	return nil
}

func (s *Service) ProductsAdd(ctx context.Context, ps []product.PreAddProduct) error {
	products := []product.Product{}
	for _, p := range ps {
		product := product.New(p)
		products = append(products, *product)
	}

	err := s.repo.AddProducts(ctx, products)
	if err != nil {
		log.Println(err)
		return ErrServer
	}
	return nil
}

func (s *Service) UsersList(ctx context.Context) ([]user.UserCard, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, ErrServer
	}

	return users, nil
}
