package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"store/product"
	"store/store"
	"store/user"
	"strconv"

	"github.com/go-chi/chi"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]user.User, error)
	AddUser(ctx context.Context, u user.User) error
	AddCard(ctx context.Context, c user.Card) error
	GetUser(ctx context.Context, fn, email, pw string) (user.User, error)
	AddProduct(ctx context.Context, p product.Product) error
	GetProduct(ctx context.Context, name string) (product.Product, error)
	ListProducts(ctx context.Context) ([]product.Product, error)
	SellProduct(ctx context.Context, sale store.Sales, product product.Product) error
}

type Handler struct {
	repo Repository
}

var neededUser = user.User{}

func (h Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalln()
	}
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	tempUser := struct {
		Name       string `json:"full_name,omitempty"`
		Email      string `json:"email,omitempty"`
		Password   string `json:"password,omitempty"`
		CardNumber int    `json:"card_number,omitempty"`
		Balance    int    `json:"balance,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalln(err)
	}

	card := user.NewCard(tempUser.CardNumber, tempUser.Balance)
	err := h.repo.AddCard(r.Context(), *card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}

	newUser := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.Email,
		*card,
	)
	err = h.repo.AddUser(r.Context(), *newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}
	neededUser = *newUser

}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	tempUser := struct {
		Name     string `json:"full_name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalln(err)
	}

	gotUser, err := h.repo.GetUser(r.Context(),
		tempUser.Name, tempUser.Email, tempUser.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}
	neededUser = gotUser
}

func (h Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	p := product.Product{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}

	product := product.New(p.Name, p.Quantity, p.Price)
	err := h.repo.AddProduct(r.Context(), *product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

func BuyProduct(repo Repository, u user.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productName := r.URL.Query().Get("name")
		quantity := r.URL.Query().Get("quantity")
		q, err := strconv.Atoi(quantity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatalln(err)
		}
		GotProduct, err := repo.GetProduct(r.Context(), productName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				p := product.Product{Name: productName}
				_, added := store.Sell(p, q, u, false, w)
				err := repo.AddProduct(r.Context(), added)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					log.Fatalln(err)
				}
			}
		}

		sales, soldProduct := store.Sell(GotProduct, q, u, true, w)
		err = repo.SellProduct(r.Context(), sales, soldProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalln(err)
		}

	}
}

func NewRouter(repo Repository) http.Handler {
	r := chi.NewRouter()
	h := Handler{repo: repo}

	r.Post("/register", h.AddUser)
	r.Post("/login", h.Login)
	r.Get("/users", h.ListUsers)
	r.Get("/buy/", BuyProduct(h.repo, neededUser))
	r.Post("/add/product", h.AddProduct)
	
	return r
}
