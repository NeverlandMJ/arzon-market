package store

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"store/product"
	"store/user"
	"time"
)

type Sales struct {
	CustomerID   string    `json:"customer_id,omitempty"`
	ProductID    string    `json:"product_id,omitempty"`
	SoldQuantity int       `json:"sold_quantity,omitempty"`
	Profit       int       `json:"profit,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}

func Sell(p product.Product, quantity int, u user.User, exist bool, w http.ResponseWriter) (Sales, product.Product) {
	profit := p.Price * quantity
	if !exist {
		io.WriteString(w, "product with this name not found\n")
		io.WriteString(w, "we will bring this product the next time")
		newProduct := product.New(p.Name, quantity, generateRandomPrice(1000, 2000))
		return Sales{}, *newProduct
	}
	if p.Quantity < quantity {
		msg := fmt.Sprintf("not enough products: only %d left\n", p.Quantity)
		io.WriteString(w, msg)
		io.WriteString(w, "we will bring this product the next time")
		newProduct := product.New(p.Name, quantity, generateRandomPrice(1000, 2000))
		return Sales{}, *newProduct
	}

	

	sales := Sales{
		CustomerID:   u.ID,
		ProductID:    p.ID,
		SoldQuantity: quantity,
		Profit:       profit,
		Time:         time.Now(),
	}

	return sales, p
}

// generates random price for the new product
func generateRandomPrice(down, up int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(up+down) - down
	fmt.Println(n)
	return n
}
