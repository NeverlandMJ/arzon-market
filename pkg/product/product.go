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
	Category      string `json:"category,omitempty"`
}

type PreAddProduct struct {
	Name        string `json:"name,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	ImageLink   string `json:"image_link,omitempty"`
	Price       int    `json:"price,omitempty"`
}

func New(p PreAddProduct) *Product {
	id := uuid.New()
	op := p.Price / 2
	return &Product{
		ID:            id.String(),
		Name:          p.Name,
		Description:   p.Description,
		Quantity:      p.Quantity,
		Price:         p.Price,
		OriginalPrice: op,
		ImageLink:     p.ImageLink,
		Category:      p.Category,
	}
}
