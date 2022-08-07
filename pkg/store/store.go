package store

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/NeverlandMJ/arzon-market/pkg/product"
)

var errNotEnoughQuantity = errors.New("we don't have enough product")

type Sales struct {
	CustomerID   string    `json:"customer_id,omitempty"`
	ProductID    string    `json:"product_id,omitempty"`
	SoldQuantity int       `json:"sold_quantity,omitempty"`
	Profit       int       `json:"profit,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}

func Sell(p product.Product, quantity int, uID string) (Sales, product.Product, error) {
	profit := p.Price * quantity

	if p.Quantity < quantity {
		return Sales{}, p, errNotEnoughQuantity
	}

	sales := Sales{
		CustomerID:   uID,
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
