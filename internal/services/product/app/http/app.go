package httpapp

import (
	"local/gorm-example/internal/lib/jwt"
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer     *gin.Engine
	address        string
	productService product.Service
	secret         string
}

type productResponse struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

type productCreateRequest struct {
	Code  string `json:"code" binding:"required"`
	Price uint   `json:"price" binding:"required"`
}

type productUpdateRequest struct {
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

type productParams struct {
	ID uint `uri:"id" binding:"required"`
}

func New(productService product.Service, address, secret string) *App {
	router := gin.Default()

	return &App{
		httpServer:     router,
		address:        address,
		productService: productService,
		secret:         secret,
	}
}

func (a *App) RegisterHandlers() {
	a.httpServer.Use(a.jwtTokenCheck)
	a.httpServer.GET("/products", a.list)
	a.httpServer.GET("/products/:id", a.findOne)
	a.httpServer.POST("/products", a.create)
	a.httpServer.DELETE("/products/:id", a.delete)
	a.httpServer.PATCH("/products/:id", a.update)
}

func (a *App) MustRun() {
	if err := a.httpServer.Run(a.address); err != nil {
		panic(err)
	}
}

func (a *App) list(c *gin.Context) {
	products, err := a.productService.List()

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err})
		return
	}

	var response []*productResponse

	for _, product := range products {
		response = append(response, toProductResponse(product))
	}

	c.JSON(http.StatusOK, response)
}

func (a *App) findOne(c *gin.Context) {
	var params productParams

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	product, err := a.productService.GetOne(params.ID)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(*product))
}

func (a *App) create(c *gin.Context) {
	var body productCreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := a.productService.Create(&models.Product{
		Price: body.Price,
		Code:  body.Code,
	}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

func (a *App) delete(c *gin.Context) {
	var params productParams

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := a.productService.Delete(params.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (a *App) update(c *gin.Context) {
	var params productParams
	var body productUpdateRequest

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := a.productService.Update(params.ID, &models.Product{
		Price: body.Price,
		Code:  body.Code,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (a *App) jwtTokenCheck(c *gin.Context) {
	if err := jwt.ValidateToken(c.Request.Header.Get("Authorization"), a.secret); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func toProductResponse(product models.Product) *productResponse {
	return &productResponse{
		ID:    product.ID,
		Price: product.Price,
		Code:  product.Code,
	}
}
