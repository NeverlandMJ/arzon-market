package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NeverlandMJ/arzon-market/pkg/middlewares"
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
}

type message struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// @title Arzon-market API
// @version 1.0
// @description online meva va poliz mahsulotlari sotiladigan magazen APIsi

// @contact.name Sunbula Hasanova
// @contact.url https://t.me/Neverland_MJ
// @contact.email khasanovasumbula@gmail.com

// @host localhost:8081
// @BasePath /api
// @query.collection.format multi
func NewRouter(serv service.Handler) *gin.Engine {
	router := gin.Default()
	s := api{serve: serv}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/auth")
	auth.POST("/register", s.SignUp)
	auth.POST("/login", s.Login)
	auth.DELETE("/logout", s.Logout)

	authored := router.Group("/api")
	authored.Use(middlewares.Authentication)
	authored.POST("/add/card", s.AddCard)
	authored.GET("/buy/", s.BuyProduct)

	router.GET("/api/product/:id", s.GetProduct)
	router.GET("/api/product/list", s.ListProducts)

	protected := router.Group("api/admin")
	protected.Use(middlewares.CheckAdmin)
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
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure		 422 {object} message
// @Failure      500  {object}  message
// @Router       /auth/register [POST]
func (a *api) SignUp(c *gin.Context) {
	tempUser := user.PreSignUpUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{
			Message: "invalid json",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	err := a.serve.CreateUser(c.Request.Context(), tempUser)

	if errors.Is(err, service.ErrUserExist) {
		r := message{
			Message: "user mavjud",
			Success: false,
		}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	} else if errors.Is(err, service.ErrInvalidUser) {
		r := message{
			Message: "to'liq ma'lumot kiritilmagan",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		return
	} else if err != nil {
		r := message{
			Message: "serverda xatolik mavjud",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r := message{
		Message: "user succesfully created",
		Success: true,
	}
	c.JSON(http.StatusCreated, r)
}

// @Summary      sign in
// @Description  user login qilishi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body user.PreLoginUser  true  "User info"
// @Success      200  {object}  message
// @Failure      400  {object} 	message
// @Failure      401  {object}  message
// @Failure      500  {object} 	message
// @Router       /auth/login [POST]
func (a *api) Login(c *gin.Context) {
	tempUser := user.PreLoginUser{}

	if err := c.BindJSON(&tempUser); err != nil {
		r := message{
			Message: "invalid json",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	token, err := a.serve.LoginUser(c.Request.Context(), tempUser)

	if errors.Is(err, service.ErrUserNotExist) {
		r := message{
			Message: "user mavjud emas",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, r)
		return
	} else if errors.Is(err, service.ErrServer) {
		r := message{
			Message: "serverda xatolik mavjud",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	c.SetCookie(
		"token",
		token,
		3600,
		"/",
		"localhost",
		false,
		true,
	)

	r := message{
		Message: "succesfully loged in",
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

// @Summary      log out
// @Description  user log out qilishi
// @Tags         auth
// @Success      200  {object}  user.User
// @Router       /auth/logout [DELETE]
func (a *api) Logout(c *gin.Context) {
	cook, _ := c.Request.Cookie("token")
	cook.Value = "deleted"
	r := message{
		Message: "loged out",
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

// @Summary      plastik karta qo'shish
// @Description  user o'zining plastik kartasini kiritishi
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body  user.PreAddCard true "Card info"
// @Success      201  {object}  message
// @Failure      400  {object}  message
// @Failure      401  {object}  message
// @Failure      500  {object}  message
// @Router       /add/card [post]
func (a *api) AddCard(c *gin.Context) {
	v, ok := c.Get("claims")
	if !ok {
		r := message{
			Message: "user token mavjud emas",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, r)
		return
	}

	claims, ok := v.(*service.Claims)
	fmt.Println(claims)
	if !ok {
		r := message{
			Message: "looks like cookie isn't set",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, r)
		return
	}

	var card user.PreAddCard

	if err := c.BindJSON(&card); err != nil {
		r := message{
			Message: "invalid json",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	_, err := a.serve.CreateCard(c.Request.Context(), claims.ID, card)

	if err != nil {
		r := message{
			Message: "error in creating new card",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r := message{
		Message: "card is added",
		Success: true,
	}
	c.JSON(http.StatusCreated, r)
}

// @Summary      buy product
// @Description  produkta sotib olish
// @Tags         user
// @Produce      json
// @Param        id query string true "product id"
// @Param        quantity query int true "product quantity"
// @Success      200  {object}  message
// @Failure      400  {object} 	message
// @Failure      401  {object}  message
// @Failure      500  {object}  message
// @Router       /buy/ [GET]
func (a *api) BuyProduct(c *gin.Context) {
	v, ok := c.Get("claims")
	if !ok {
		return
	}

	claims, ok := v.(*service.Claims)
	// fmt.Println(claims)
	if !ok {
		r := message{
			Message: "looks like cookie isn't set",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, r)
		return
	}

	productID, ok := c.GetQuery("id")
	if !ok {
		r := message{
			Message: "empty query",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	quantity, ok := c.GetQuery("quantity")
	if !ok {
		r := message{
			Message: "empty query",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	q, err := strconv.Atoi(quantity)
	if err != nil {
		r := message{
			Message: "invalid query",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		// fmt.Println(err)
		return
	}

	err = a.serve.SellProduct(c.Request.Context(), productID, q, claims.ID)

	if err != nil {
		if errors.Is(err, service.ErrProductNotExist) {
			r := message{
				Message: "product mavjud emas",
				Success: false,
			}
			c.JSON(http.StatusBadRequest, r)
			return
		} else if errors.Is(err, service.ErrQuantityExceeded) {
			r := message{
				Message: "product miqdori bazada yetarli emas",
				Success: false,
			}
			c.JSON(http.StatusBadRequest, r)
			return
		} else if errors.Is(err, service.ErrServer) {
			r := message{
				Message: "server xatoligi",
				Success: false,
			}
			c.JSON(http.StatusInternalServerError, r)
			return
		} else if errors.Is(err, service.ErrNotEnoughBalance) {
			r := message{
				Message: "userni puli yetmaydi",
				Success: false,
			}
			c.JSON(http.StatusBadRequest, r)
			return
		} else if errors.Is(err, service.ErrCardNotExist) {
			r := message{
				Message: "userda karta mavjud emas",
				Success: false,
			}
			c.JSON(http.StatusBadRequest, r)
			return
		}
	}

	r := message{
		Message: "product sotildi",
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

// @Summary      hamma produktalarni listi
// @Description  barcha produktalarni ko'rsatish
// @Tags         public
// @Produce      json
// @Success      200  {object}  []product.Product
// @Failure      500  {object}  message
// @Router       /product/list [GET]
func (a *api) ListProducts(c *gin.Context) {

	products, err := a.serve.AllProducts(c.Request.Context())
	if err != nil {
		r := message{
			Message: "hamma productlar haqida ma'lumot chiqmadi",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	c.JSON(http.StatusOK, products)

}

// @Summary      produkta
// @Description  bitta produkta haqida ma'lumotlarni olish
// @Tags         public
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
		r := message{
			Message: "invalid params",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
	}

	p, err := a.serve.GetOneProductInfo(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotExist) {
			r := message{
				Message: "product mavjud emas",
				Success: false,
			}
			c.JSON(http.StatusBadGateway, r)
			fmt.Println(err)
			return
		} else {
			r := message{
				Message: "server xatoligi",
				Success: false,
			}
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
// @Param        request body product.PreAddProduct true "Product info"
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /admin/add/product [POST]
func (a *api) AddProduct(c *gin.Context) {
	_, ok := c.Get("claims")
	if !ok {
		r := message{
			Message: "user admin emas",
			Success: false,
		}
		c.JSON(http.StatusMethodNotAllowed, r)
		return

	}

	p := product.PreAddProduct{}

	if err := c.BindJSON(&p); err != nil {
		r := message{
			Message: "invalid json",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}

	err := a.serve.ProductAdd(c.Copy().Request.Context(), p)
	if err != nil {
		r := message{
			Message: "server xatoligi",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{
		Message: "product qo'shildi",
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

// @Summary      produktalar qo'shish
// @Description  bir nechta produktalarni qo'shish
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body []product.PreAddProduct true "Product info"
// @Success      200  {object}  message
// @Failure      400  {object}  message
// @Failure      500  {object}  message
// @Router       /admin/add/list/product [POST]
func (a *api) AddProducts(c *gin.Context) {
	_, ok := c.Get("claims")
	if !ok {
		r := message{
			Message: "user admin emas",
			Success: false,
		}
		c.JSON(http.StatusMethodNotAllowed, r)
		return

	}
	tempProducts := []product.PreAddProduct{}

	if err := c.BindJSON(&tempProducts); err != nil {
		r := message{
			Message: "invalid json",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, r)
		fmt.Println(err)
		return
	}
	err := a.serve.ProductsAdd(c.Request.Context(), tempProducts)
	if err != nil {
		r := message{
			Message: "server xatoligi",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		fmt.Println(err)
		return
	}

	r := message{
		Message: "productlar qo'shildi",
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

// @Summary      hamma userlar ro'yxati
// @Description  hamma userlar ro'yxatini chiqarish
// @Tags         admin
// @Produce      json
// @Success      200  {object}  []user.UserCard
// @Failure      405  {object}  message
// @Failure      500  {object}  message
// @Router       /admin/users [GET]
func (a *api) ListUsers(c *gin.Context) {
	_, ok := c.Get("claims")
	if !ok {
		return
	}
	users, err := a.serve.UsersList(c.Request.Context())
	if err != nil {
		r := message{
			Message: "server xatoligi",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	c.JSON(http.StatusOK, users)
}
