package product

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Quantity      int    `json:"quantity,omitempty"`
	Price         int    `json:"price,omitempty"`
	OriginalPrice int    `json:"original_price,omitempty"`
}

func New(name string, q, p int) *Product {
	id := uuid.New()
	op := p/2
	return &Product{
		ID:            id.String(),
		Name:          name,
		Quantity:      q,
		Price:         p,
		OriginalPrice: op,
	}
}

// generates random price for the new product
func generateRandomPrice(down, up int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(up+down) - down

	return n
}
