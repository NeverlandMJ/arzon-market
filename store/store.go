package store

import (
	"errors"
	"fmt"
	"math/rand"
	"store/product"
	"store/user"
	"time"
)

var errNotEnoughQuantity = errors.New("we don't have enough product")

type Sales struct {
	CustomerID   string    `json:"customer_id,omitempty"`
	ProductID    string    `json:"product_id,omitempty"`
	SoldQuantity int       `json:"sold_quantity,omitempty"`
	Profit       int       `json:"profit,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}

func Sell(p product.Product, quantity int, u user.User) (Sales, product.Product, error) {
	profit := p.Price * quantity

	if p.Quantity < quantity {
		return Sales{}, p, errNotEnoughQuantity
	}

	sales := Sales{
		CustomerID:   u.ID,
		ProductID:    p.ID,
		SoldQuantity: quantity,
		Profit:       profit,
		Time:         time.Now(),
	}

	return sales, p, nil
}

// generates random price for the new product
func generateRandomPrice(down, up int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(up+down) - down
	fmt.Println(n)
	return n
}
