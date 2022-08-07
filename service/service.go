package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/NeverlandMJ/arzon-market/pkg/product"
	"github.com/NeverlandMJ/arzon-market/pkg/store"
	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/NeverlandMJ/arzon-market/storage"
)

var ErrUserExist = errors.New("user already exist")
var ErrInvalidUser = errors.New("empty faild")
var ErrUserNotExist = errors.New("user doesn't exist")
var ErrServer = errors.New("server error")
var ErrProductNotExist = errors.New("product doesn't exist")
var ErrQuantityExceeded = errors.New("quantity exceeded")

type Handler interface {
	CreateUser(ctx context.Context, tempUser user.PreSignUpUser) (user.User, error)
	LoginUser(ctx context.Context, tempUser user.PreLoginUser) (user.User, error)
	CreateCard(ctx context.Context, owner user.User, card user.Card) (user.Card, error)
	SellProduct(ctx context.Context, productName string, quantity int, client user.User) error
	AllProducts(ctx context.Context) ([]product.Product, error)
	GetOneProductInfo(ctx context.Context, name string) (product.Product, error)
	ProductAdd(ctx context.Context, p product.Product) error
	ProductsAdd(ctx context.Context, ps []product.Product) error
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

func (s *Service) CreateUser(ctx context.Context, tempUser user.PreSignUpUser) (user.User, error) {
	fmt.Println(tempUser)
	if tempUser.Name == "" || tempUser.PhoneNumber == "" || tempUser.Password == "" {
		return user.User{}, ErrInvalidUser
	}

	newUser := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.PhoneNumber,
	)
	err := s.repo.AddUser(ctx, *newUser)

	if err != nil {
		log.Println("Service(): ", err)
		return user.User{}, ErrUserExist
	}

	return *newUser, nil
}

func (s *Service) LoginUser(ctx context.Context, tempUser user.PreLoginUser) (user.User, error) {
	if tempUser.PhoneNumber == "" || tempUser.Password == "" {
		return user.User{}, ErrInvalidUser
	}

	u, err := s.repo.GetUser(ctx, tempUser.PhoneNumber, tempUser.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, ErrUserNotExist
		} else {
			log.Println("Service(): ", err)
			return user.User{}, ErrServer
		}
	}

	return u, nil
}

func (s *Service) CreateCard(ctx context.Context, owner user.User, card user.Card) (user.Card, error) {
	card.OwnerID = owner.ID
	newCard := user.NewCard(card.CardNumber, card.Balance, card.OwnerID)

	err := s.repo.AddCard(ctx, *newCard)
	if err != nil {
		log.Println("Service(): ", err)
		return user.Card{}, ErrServer
	}

	return *newCard, nil
}

func (s *Service) SellProduct(ctx context.Context, productName string, quantity int, client user.User) error {
	gotProduct, err := s.repo.GetProduct(ctx, productName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotExist
		} else {
			log.Println(err)
			return ErrServer
		}
	}

	sales, soldProduct, err := store.Sell(gotProduct, quantity, client)

	if err != nil {
		return ErrQuantityExceeded
	}

	err = s.repo.SellProduct(ctx, sales, soldProduct)
	if err != nil {
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

func (s *Service) ProductAdd(ctx context.Context, p product.Product) error {
	product := product.New(p.Name, p.Description, p.ImageLink, p.Category, p.Quantity, p.Price)
	err := s.repo.AddProduct(ctx, *product)
	if err != nil {
		log.Println(err)
		return ErrServer
	}
	return nil
}

func (s *Service) ProductsAdd(ctx context.Context, ps []product.Product) error {
	products := []product.Product{}
	for _, p := range ps {
		product := product.New(p.Name, p.Description, p.ImageLink, p.Category, p.Quantity, p.Price)
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
