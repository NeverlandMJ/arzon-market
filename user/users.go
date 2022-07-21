package user

import "github.com/google/uuid"

type User struct {
	ID       string `json:"id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	CardID   string `json:"card_id,omitempty"`
	IsAdmin  bool `json:"is_admin,omitempty"`
}

type Card struct {
	ID         string `json:"id,omitempty"`
	CardNumber int    `json:"card_number,omitempty"`
	Balance    int    `json:"balance,omitempty"`
}

type SellProductRequest struct {
	ProductName string
	Quantity    int
}

type Cart []SellProductRequest

// AddToCart adds bought products into client's cart
func (c *Cart) AddToCart(r SellProductRequest) {
	*c = append(*c, r)
}


func NewCard(cn, bln int) *Card {
	id := uuid.New()
	return &Card{
		ID: id.String(),
		CardNumber: cn,
		Balance: bln,
	}
}

func NewUser(fn, pw, email string, NewCard Card) *User {
	id := uuid.New()
	return &User{
		ID: id.String(),
		FullName: fn,
		Password: pw,
		Email: email,
		CardID: NewCard.ID,
		IsAdmin: false,
	}
}

