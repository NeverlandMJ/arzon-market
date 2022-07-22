package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"store/product"
	"store/store"
	"store/user"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]user.UserCard, error)
	AddUser(ctx context.Context, u user.User, c user.Card) error
	AddCard(ctx context.Context, c user.Card) error
	GetUser(ctx context.Context, fn, email, pw string) (user.User, error)
	AddProduct(ctx context.Context, p product.Product) error
	GetProduct(ctx context.Context, name string) (product.Product, error)
	ListProducts(ctx context.Context) ([]product.Product, error)
	SellProduct(ctx context.Context, sale store.Sales, product product.Product) error
	AddProducts(ctx context.Context, ps []product.Product) error
}

type Handler struct {
	repo Repository
	user user.User
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	tempUser := struct {
		Name       string `json:"full_name,omitempty"`
		Email      string `json:"email,omitempty"`
		Password   string `json:"password,omitempty"`
		CardNumber string `json:"card_number,omitempty"`
		Balance    int    `json:"balance,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	card := user.NewCard(tempUser.CardNumber, tempUser.Balance)
	newUser := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.Email,
		*card,
	)
	err := h.repo.AddUser(r.Context(), *newUser, *card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	h.user = *newUser
	io.WriteString(w, "userregistreted")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	tempUser := struct {
		Name     string `json:"full_name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	gotUser, err := h.repo.GetUser(r.Context(),
		tempUser.Name, tempUser.Email, tempUser.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	h.user = gotUser
	io.WriteString(w, "loggedin")
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	p := product.Product{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	product := product.New(p.Name, p.Description, p.Quantity, p.Price)
	err := h.repo.AddProduct(r.Context(), *product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	io.WriteString(w, "productadded")
}

func (h *Handler) AddProducts(w http.ResponseWriter, r *http.Request) {
	tempProducts := []product.Product{}
	if err := json.NewDecoder(r.Body).Decode(&tempProducts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	products := []product.Product{}
	for _, p := range tempProducts {
		product := product.New(p.Name, p.Description, p.Quantity, p.Price)
		products = append(products, *product)
	}

	h.repo.AddProducts(r.Context(), products)

}

func (h *Handler) BuyProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	productName := r.URL.Query().Get("name")
	quantity := r.URL.Query().Get("quantity")
	q, err := strconv.Atoi(quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	GotProduct, err := h.repo.GetProduct(r.Context(), productName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p := product.Product{Name: productName}
			_, added := store.Sell(p, q, h.user, false, w)
			err := h.repo.AddProduct(r.Context(), added)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				panic(err)
			}
		}
	}

	sales, soldProduct := store.Sell(GotProduct, q, h.user, true, w)

	err = h.repo.SellProduct(r.Context(), sales, soldProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	io.WriteString(w, "productissold")

}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	products, err := h.repo.ListProducts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	productName := r.URL.Query().Get("name")
	p, err := h.repo.GetProduct(r.Context(), productName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func NewRouter(repo Repository) http.Handler {
	r := chi.NewRouter()
	h := Handler{repo: repo}
	r.Use(middleware.Logger)
	r.Post("/register", h.AddUser)
	r.Post("/login", h.Login)
	r.Get("/users", h.ListUsers)
	r.Get("/buy/", h.BuyProduct)
	r.Get("/product/", h.GetProduct)
	r.Post("/product/add", h.AddProduct)
	r.Post("/product/add/list", h.AddProducts)
	r.Get("/product/list", h.ListProducts)

	return r
}
