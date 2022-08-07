package storage

import (
	"context"

	"github.com/NeverlandMJ/arzon-market/pkg/product"
	"github.com/NeverlandMJ/arzon-market/pkg/store"
	"github.com/NeverlandMJ/arzon-market/pkg/user"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]user.UserCard, error)
	AddUser(ctx context.Context, u user.User) error
	GetUser(ctx context.Context, email, pw string) (user.User, error)
	AddProduct(ctx context.Context, p product.Product) error
	GetProduct(ctx context.Context, name string) (product.Product, error)
	ListProducts(ctx context.Context) ([]product.Product, error)
	SellProduct(ctx context.Context, sale store.Sales, product product.Product) error
	AddProducts(ctx context.Context, ps []product.Product) error
	AddCard(ctx context.Context, c user.Card) error
	GetUserByID(id string) (user.User, error)
}
