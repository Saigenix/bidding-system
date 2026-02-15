package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saigenix/bidding-system/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required" example:"Gaming Laptop"`
	Description string `json:"description" example:"High-performance gaming laptop"`
}

// Create godoc
// @Summary      Create a product
// @Description  Create a new product that can be listed for auction
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request  body      CreateProductRequest  true  "Product details"
// @Success      201      {object}  domain.Product
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	product, err := h.productService.CreateProduct(c.Request.Context(), req.Name, req.Description, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Get godoc
// @Summary      Get a product
// @Description  Get a product by its ID
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// List godoc
// @Summary      List all products
// @Description  Get a list of all available products
// @Tags         Products
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.productService.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
