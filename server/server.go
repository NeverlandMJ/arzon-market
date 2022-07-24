package server

import (
	"context"
	"encoding/json"
	"fmt"
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
	AddUser(ctx context.Context, u user.User) error
	GetUser(ctx context.Context, email, pw string) (user.User, error)
	AddProduct(ctx context.Context, p product.Product) error
	GetProduct(ctx context.Context, name string) (product.Product, error)
	ListProducts(ctx context.Context) ([]product.Product, error)
	SellProduct(ctx context.Context, sale store.Sales, product product.Product) error
	AddProducts(ctx context.Context, ps []product.Product) error
	AddCard(ctx context.Context, c user.Card) error
}

type Handler struct {
	repo Repository
	user user.User
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.user.IsAdmin{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	tempUser := struct {
		Name     string `json:"full_name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
		// CardNumber string `json:"card_number,omitempty"`
		// Balance    int    `json:"balance,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// card := user.NewCard(tempUser.CardNumber, tempUser.Balance)
	newUser := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.Email,
	)
	err := h.repo.AddUser(r.Context(), *newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	h.user = *newUser
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "userregistreted")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	tempUser := struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	gotUser, err := h.repo.GetUser(r.Context(),
		tempUser.Email, tempUser.Password)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "userdoesntexist")
		fmt.Println(err)
		return
	}
	h.user = gotUser
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "loggedin")
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	if !h.user.IsAdmin{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	p := product.Product{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	product := product.New(p.Name, p.Description, p.Quantity, p.Price)
	err := h.repo.AddProduct(r.Context(), *product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "productadded")
}

func (h *Handler) AddProducts(w http.ResponseWriter, r *http.Request) {
	if !h.user.IsAdmin{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tempProducts := []product.Product{}
	if err := json.NewDecoder(r.Body).Decode(&tempProducts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	products := []product.Product{}
	for _, p := range tempProducts {
		product := product.New(p.Name, p.Description, p.Quantity, p.Price)
		products = append(products, *product)
	}

	h.repo.AddProducts(r.Context(), products)

}

func (h *Handler) BuyProduct(w http.ResponseWriter, r *http.Request) {

	productName := r.URL.Query().Get("name")
	quantity := r.URL.Query().Get("quantity")
	q, err := strconv.Atoi(quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	GotProduct, err := h.repo.GetProduct(r.Context(), productName)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			return
		}
	}

	sales, soldProduct, err := store.Sell(GotProduct, q, h.user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "quantity exceeded")
		fmt.Println(err)
		return
	}

	err = h.repo.SellProduct(r.Context(), sales, soldProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "productissold")

}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.repo.ListProducts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productName := r.URL.Query().Get("name")
	p, err := h.repo.GetProduct(r.Context(), productName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request)  {
	var c user.Card
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	c.OwnerID = h.user.ID
	newCard := user.NewCard(c.CardNumber, c.Balance, c.OwnerID)

	err := h.repo.AddCard(r.Context(), *newCard)
	if err != nil {
		fmt.Println(err)
		return
	}	
	w.WriteHeader(http.StatusCreated)
}

func NewRouter(repo Repository) http.Handler {
	r := chi.NewRouter()
	h := Handler{repo: repo}
	r.Use(middleware.Logger)
	r.Post("/register", h.SignUp)//
	r.Post("/login", h.Login)//
	r.Post("/add/card", h.AddCard)//
	r.Get("/users", h.ListUsers)
	r.Get("/buy/", h.BuyProduct)
	r.Get("/product/", h.GetProduct)//
	r.Post("/add/product", h.AddProduct)//
	r.Post("/add/list/product", h.AddProducts)//
	r.Get("/product/list", h.ListProducts)//

	return r
}
