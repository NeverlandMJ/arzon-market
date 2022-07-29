package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NeverlandMJ/arzon-market/product"
	"github.com/NeverlandMJ/arzon-market/store"
	"github.com/NeverlandMJ/arzon-market/user"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/NeverlandMJ/arzon-market/docs"

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

// message represents request response with a message
type message struct {
	Message string `json:"message"`
}

// @title Arzon-market API
// @version 1.0
// @description online meva va poliz mahsulotlari sotiladigan magazen APIsi

// @contact.name Sunbula Hasanova
// @contact.url https://t.me/Neverland_MJ
// @contact.email khasanovasumbula@gmail.com

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func NewRouter(repo Repository) *gin.Engine {
	router := gin.Default()
	h := Handler{repo: repo}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register", h.SignUp)  
	router.POST("/login", h.Login)      
	router.POST("/add/card", h.AddCard) 
	router.GET("/buy/", h.BuyProduct)
	router.GET("/product/:id", h.GetProduct)           
	router.GET("/product/list", h.ListProducts)     
	
	router.GET("/users", h.ListUsers)
	router.POST("/add/product", h.AddProduct)       
	router.POST("/add/list/product", h.AddProducts) 

	return router
}


// @Summary      hamma userlar ro'yxati
// @Description  hamma userlar ro'yxatini chiqarish
// @Tags         admin
// @Produce      json
// @Success      200  {object}  []user.UserCard
// @Failure      405  {object}  message
// @Failure      500  {object}  message
// @Router       /users [GET]
func (h *Handler) ListUsers(c *gin.Context) {
	if !h.user.IsAdmin {
		r := message{"method is only allowed to admins"}
		c.JSON(http.StatusMethodNotAllowed, r)
		return
	}

	users, err := h.repo.ListUsers(c.Request.Context())
	if err != nil {
		r := message{"error while listing users"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, users)
}


// @Summary      sign up
// @Description  user registratsiyadan o'tishi
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body user.PreSignUpUser  true  "User info"
// @Success      200  {object}  user.PreSignUpUser
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /register [POST]
func (h *Handler) SignUp(c *gin.Context) {
	tempUser := user.PreSignUpUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	newUser := user.NewUser(
		tempUser.Name,
		tempUser.Password,
		tempUser.Email,
	)
	err := h.repo.AddUser(c.Request.Context(), *newUser)
	if err != nil {
		r := message{"error in adding user to database"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}
	h.user = *newUser

	// r := message{"user registrated"}
	c.JSON(http.StatusCreated, *newUser)
}


// @Summary      sign in
// @Description  user login qilishi
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body user.PreLoginUser  true  "User info"
// @Success      200  {object}  user.PreLoginUser
// @Failure      400  {object} 	message
// @Failure      404  {object}  message
// @Failure      500  {object} 	message
// @Router       /login [POST]
func (h *Handler) Login(c *gin.Context) {
	tempUser := user.PreLoginUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	gotUser, err := h.repo.GetUser(c.Request.Context(), tempUser.Email, tempUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			r := message{"user doesn't exist"}
			c.JSON(http.StatusBadRequest, r)
			fmt.Println(err)
			return
		}else {
			r := message{err.Error()}
			c.JSON(http.StatusInternalServerError, r)
			fmt.Println(err)
			return
		}
		
	}
	h.user = gotUser

	// r := message{"logged in"}
	c.JSON(http.StatusOK, gotUser)
}


// @Summary      plastik karta qo'shish
// @Description  user o'zining plastik kartasini kiritishi
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body  user.Card true "Card info"
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /add/card [get]
func (h *Handler) AddCard(c *gin.Context) {
	var card user.Card

	if err := c.BindJSON(&card); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	card.OwnerID = h.user.ID
	newCard := user.NewCard(card.CardNumber, card.Balance, card.OwnerID)

	err := h.repo.AddCard(c.Request.Context(), *newCard)
	if err != nil {
		r := message{"error in creating new card"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"card is added"}
	c.JSON(http.StatusCreated, r)
}

// @Summary      produkta qo'shish
// @Description  bitta produkta qo'shish
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body product.Product true "Product info"
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /add/product [POST]
func (h *Handler) AddProduct(c *gin.Context) {
	if !h.user.IsAdmin {
		r := message{"method is only allowed to admins"}
		c.JSON(http.StatusMethodNotAllowed, r)
		return
	}
	p := product.Product{}

	if err := c.BindJSON(&p); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	product := product.New(p.Name, p.Description, p.ImageLink, p.Quantity, p.Price)
	err := h.repo.AddProduct(c.Request.Context(), *product)
	if err != nil {
		r := message{"error in creating a new food"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"product added"}
	c.JSON(http.StatusOK, r)
}

// @Summary      produktalar qo'shish
// @Description  bir nechta produktalarni qo'shish
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body []product.Product true "Product info"
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /add/list/product [POST]
func (h *Handler) AddProducts(c *gin.Context) {
	if !h.user.IsAdmin {
		r := message{"method is only allowed to admins"}
		c.JSON(http.StatusMethodNotAllowed, r)
		return
	}
	tempProducts := []product.Product{}

	if err := c.BindJSON(&tempProducts); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	products := []product.Product{}
	for _, p := range tempProducts {
		product := product.New(p.Name, p.Description, p.ImageLink, p.Quantity, p.Price)
		products = append(products, *product)
	}

	err := h.repo.AddProducts(c.Request.Context(), products)
	if err != nil {
		r := message{"error in adding new foods"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"products added"}
	c.JSON(http.StatusOK, r)
}


// @Summary      buy product
// @Description  produkta sotib olish
// @Tags         user
// @Produce      json
// @Param        name query string quantity query int false "Buy product"
// @Success      200  {object}  message
// @Failure      400  {object} 	message
// @Failure      404  {object}  message
// @Failure      500  {object}  message
// @Router       /buy/ [GET]
func (h *Handler) BuyProduct(c *gin.Context) {
	productName, ok := c.GetQuery("name")
	if !ok {
		r := message{"empty query"}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	quantity, ok := c.GetQuery("quantity")
	if !ok {
		r := message{"empty query"}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	q, err := strconv.Atoi(quantity)
	if err != nil {
		r := message{"invalid query"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	GotProduct, err := h.repo.GetProduct(c.Request.Context(), productName)
	if err != nil {
		if err == sql.ErrNoRows {
			r := message{"no such product"}
			c.JSON(http.StatusNotExtended, r)
			fmt.Println(err)
			return
		} else {
			r := message{"error fetching data"}
			c.JSON(http.StatusInternalServerError, r)
			fmt.Println(err)
			return
		}
	}

	sales, soldProduct, err := store.Sell(GotProduct, q, h.user)
	if err != nil {
		r := message{"quantity exceeded"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	err = h.repo.SellProduct(c.Request.Context(), sales, soldProduct)
	if err != nil {
		r := message{"selling data error"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"product sold"}
	c.JSON(http.StatusOK, r)
}


// @Summary      hamma produktalarni listi
// @Description  barcha produktalarni ko'rsatish
// @Tags         user
// @Produce      json
// @Success      200  {object}  []product.Product
// @Failure      500  {object}  message
// @Router       /product/list [GET]
func (h *Handler) ListProducts(c *gin.Context) {

	products, err := h.repo.ListProducts(c.Request.Context())
	if err != nil {
		r := message{"internal error getteing list of products"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, products)

}

// @Summary      produkta 
// @Description  bitta produkta haqida ma'lumotlarni olish
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  product.Product
// @Failure      400  {object}  message
// @Failure      404  {object}  message
// @Failure      500  {object}  message
// @Router       /product/{id} [GET]
func (h *Handler) GetProduct(c *gin.Context) {
	productID, ok := c.Params.Get("name")
	if !ok {
		r := message{"invalid params"}
		c.JSON(http.StatusBadRequest, r)
	}

	p, err := h.repo.GetProduct(c.Request.Context(), productID)
	if err != nil {
		if err == sql.ErrNoRows {
			r := message{"no such product"}
			c.JSON(http.StatusNotExtended, r)
			fmt.Println(err)
			return
		}else {
			r := message{"error in fetching data"}
			c.JSON(http.StatusInternalServerError, r)
			fmt.Println(err)
			return
		}
	}

	c.JSON(http.StatusOK, p)
}



