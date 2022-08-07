package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NeverlandMJ/arzon-market/pkg/product"
	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/NeverlandMJ/arzon-market/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/NeverlandMJ/arzon-market/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type api struct {
	serve service.Handler
	user  user.User
}

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
// @BasePath /api
// @query.collection.format multi
func NewRouter(serv service.Handler) *gin.Engine {
	router := gin.Default()
	s := api{serve: serv}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api")

	public.POST("/register", s.SignUp)
	public.POST("/login", s.Login)
	public.POST("/add/card", s.AddCard)
	public.GET("//product/list", s.ListProducts)
	public.GET("/buy/", s.BuyProduct)
	public.GET("/product/:id", s.GetProduct)

	protected := router.Group("api/admin")

	protected.GET("/users", s.ListUsers)
	protected.POST("/add/product", s.AddProduct)
	protected.POST("/add/list/product", s.AddProducts)

	return router
}

// @Summary      sign up
// @Description  user registratsiyadan o'tishi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body user.PreSignUpUser  true  "User info"
// @Success      200  {object}  user.User
// @Failure      400  {object}  message
// @Failure		 422 {object} message
// @Failure      500  {object}  message
// @Router       /register [POST]
func (a *api) SignUp(c *gin.Context) {
	tempUser := user.PreSignUpUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	newUser, err := a.serve.CreateUser(c.Request.Context(), tempUser)

	if errors.Is(err, service.ErrUserExist) {
		r := message{"user mavjud"}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	} else if errors.Is(err, service.ErrInvalidUser) {
		r := message{"to'liq ma'lumot kiritilmagan"}
		c.JSON(http.StatusBadRequest, r)
		return
	} else if err != nil {
		r := message{"serverda xatolik mavjud"}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	a.user = newUser
	c.JSON(http.StatusCreated, newUser)
}

// @Summary      sign in
// @Description  user login qilishi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body user.PreLoginUser  true  "User info"
// @Success      200  {object}  user.User
// @Failure      400  {object} 	message
// @Failure      422  {object}  message
// @Failure      500  {object} 	message
// @Router       /login [POST]
func (a *api) Login(c *gin.Context) {
	tempUser := user.PreLoginUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	u, err := a.serve.LoginUser(c.Request.Context(), tempUser)

	if errors.Is(err, service.ErrUserNotExist) {
		r := message{"user mavjud emas"}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	} else if errors.Is(err, service.ErrServer) {
		r := message{"serverda xatolik mavjud"}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	a.user = u
	c.JSON(http.StatusOK, u)
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
// @Router       /add/card [post]
func (a *api) AddCard(c *gin.Context) {
	var card user.Card

	if err := c.BindJSON(&card); err != nil {
		r := message{"invalid json"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	_, err := a.serve.CreateCard(c.Request.Context(), a.user, card)

	if err != nil {
		r := message{"error in creating new card"}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r := message{"card is added"}
	c.JSON(http.StatusCreated, r)
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
func (a *api) BuyProduct(c *gin.Context) {
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

	err = a.serve.SellProduct(c.Request.Context(), productName, q, a.user)

	if err != nil {
		if errors.Is(err, service.ErrProductNotExist) {
			r := message{"product mavjud emas"}
			c.JSON(http.StatusBadRequest, r)
			return
		} else if errors.Is(err, service.ErrQuantityExceeded) {
			r := message{"product miqdori bazada yetarli emas"}
			c.JSON(http.StatusBadRequest, r)
			return
		} else if errors.Is(err, service.ErrServer) {
			r := message{"server xatoligi"}
			c.JSON(http.StatusInternalServerError, r)
			return
		}
	}

	r := message{"product sotildi"}
	c.JSON(http.StatusOK, r)
}

// @Summary      hamma produktalarni listi
// @Description  barcha produktalarni ko'rsatish
// @Tags         user
// @Produce      json
// @Success      200  {object}  []product.Product
// @Failure      500  {object}  message
// @Router       /product/list [GET]
func (a *api) ListProducts(c *gin.Context) {

	products, err := a.serve.AllProducts(c.Request.Context())
	if err != nil {
		r := message{"hamma productlar haqida ma'lumot chiqmadi"}
		c.JSON(http.StatusInternalServerError, r)
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
func (a *api) GetProduct(c *gin.Context) {
	productID, ok := c.Params.Get("id")
	if !ok {
		r := message{"invalid params"}
		c.JSON(http.StatusBadRequest, r)
	}

	p, err := a.serve.GetOneProductInfo(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotExist) {
			r := message{"product mavjud emas"}
			c.JSON(http.StatusBadGateway, r)
			fmt.Println(err)
			return
		} else {
			r := message{"server xatoligi"}
			c.JSON(http.StatusInternalServerError, r)
			fmt.Println(err)
			return
		}
	}

	c.JSON(http.StatusOK, p)
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
func (a *api) AddProduct(c *gin.Context) {
	if !a.user.IsAdmin {
		r := message{"user admin emas"}
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

	err := a.serve.ProductAdd(c.Copy().Request.Context(), p)
	if err != nil {
		r := message{"server xatoligi"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"product qo'shildi"}
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
func (a *api) AddProducts(c *gin.Context) {
	if !a.user.IsAdmin {
		r := message{"user admin emas"}
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
	err := a.serve.ProductsAdd(c.Request.Context(), tempProducts)
	if err != nil {
		r := message{"server xatoligi"}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{"productlar qo'shildi"}
	c.JSON(http.StatusOK, r)
}

// @Summary      hamma userlar ro'yxati
// @Description  hamma userlar ro'yxatini chiqarish
// @Tags         admin
// @Produce      json
// @Success      200  {object}  []user.UserCard
// @Failure      405  {object}  message
// @Failure      500  {object}  message
// @Router       /users [GET]
func (a *api) ListUsers(c *gin.Context) {
	if !a.user.IsAdmin {
		r := message{"user admin emas"}
		c.JSON(http.StatusMethodNotAllowed, r)
		return
	}

	users, err := a.serve.UsersList(c.Request.Context())
	if err != nil {
		r := message{"server xatoligi"}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	c.JSON(http.StatusOK, users)
}
