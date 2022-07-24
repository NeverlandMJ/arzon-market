package product

import (
	"github.com/google/uuid"
)

type Product struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Quantity      int    `json:"quantity,omitempty"`
	Price         int    `json:"price,omitempty"`
	OriginalPrice int    `json:"original_price,omitempty"`
	ImageLink     string `json:"image_link,omitempty"`
}

func New(name, description, link string, q, p int) *Product {
	id := uuid.New()
	op := p / 2
	return &Product{
		ID:            id.String(),
		Name:          name,
		Description:   description,
		Quantity:      q,
		Price:         p,
		OriginalPrice: op,
		ImageLink: link,
	}
}
