package user

import "github.com/google/uuid"

type User struct {
	ID       string `json:"id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty"`
}

type PreSignUpUser struct {
	Name     string `json:"full_name,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Password string `json:"password,omitempty"`
}

type PreLoginUser struct {
	PhoneNumber    string `json:"phone_number,omitempty"`
	Password string `json:"password,omitempty"`
}

type Card struct {
	ID         string `json:"id,omitempty"`
	CardNumber string `json:"card_number,omitempty"`
	Balance    int    `json:"balance,omitempty"`
	OwnerID    string `json:"owner_id,omitempty"`
}

type UserCard struct {
	FullName   string `json:"full_name,omitempty"`
	Password   string `json:"password,omitempty"`
	PhoneNumber      string `json:"email,omitempty"`
	CardNumber string `json:"card_number,omitempty"`
	Balance    int    `json:"balance,omitempty"`
}

func NewCard(cn string, bln int, userID string) *Card {
	id := uuid.New()
	return &Card{
		ID:         id.String(),
		CardNumber: cn,
		Balance:    bln,
		OwnerID:    userID,
	}
}

func NewUser(fn, pw, phone string) *User {
	id := uuid.New()
	return &User{
		ID:       id.String(),
		FullName: fn,
		Password: pw,
		PhoneNumber:    phone,
		IsAdmin:  false,
	}
}
