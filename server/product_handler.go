package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/techagentng/ecommerce-api/models"
	"github.com/techagentng/ecommerce-api/server/response"
	"gorm.io/gorm"
)

// handleCreateProduct handles the creation of a new product
// @Summary Create a new product
// @Description Only admin can create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product information"
// @Success 201 {object} Product
// @Failure 403 {object} response.ErrorResponse "Only admin users can create products"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /products [post]
func (s *Server) handleCreateProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        if userRole != "Admin" {
            response.JSON(c, "Only admin users can create products", http.StatusForbidden, nil, nil)
            return
        }

        var product models.Product
        if err := c.ShouldBindJSON(&product); err != nil {
            response.JSON(c, "Invalid JSON format", http.StatusBadRequest, nil, err)
            return
        }

        createdProduct, err := s.ProductRepo.CreateProduct(&product) 
        if err != nil {
            response.JSON(c, "Failed to create product", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Product created successfully", http.StatusCreated, createdProduct, nil)
    }
}

// handleReadProduct retrieves a product by ID
// @Summary Retrieve a product by ID
// @Description Get product details by ID (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} Product "Success"
// @Failure 403 {object} response.ErrorResponse "Only admin users can access this endpoint"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /products/{product_id} [get]
func (s *Server) handleReadProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        if userRole != "Admin" {
            response.JSON(c, "Only admin users can access this endpoint", http.StatusForbidden, nil, nil)
            return
        }

        productIDStr := c.Param("product_id")
        if productIDStr == "" {
            response.JSON(c, "Product ID cannot be empty", http.StatusBadRequest, nil, nil)
            return
        }

        productID64, err := strconv.ParseUint(productIDStr, 10, 32)
        if err != nil {
            response.JSON(c, "Invalid product ID", http.StatusBadRequest, nil, err)
            return
        }
        
        productID := uint(productID64)
        product, err := s.ProductRepo.FindProductByID(productID)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                response.JSON(c, "Product not found", http.StatusNotFound, nil, nil)
                return
            }
            response.JSON(c, "Failed to retrieve product", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Product retrieved successfully", http.StatusOK, product, nil)
    }
}

// handleUpdateProduct updates a product by ID
// @Summary Update a product by ID
// @Description Update product details by ID (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param product body Product true "Product details"
// @Success 200 {object} response.SuccessResponse "Success"
// @Failure 400 {object} response.ErrorResponse "Invalid request format"
// @Failure 403 {object} response.ErrorResponse "Only admin users can access this endpoint"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /products/{product_id} [put]
func (s *Server) handleUpdateProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        if userRole != "Admin" {
            response.JSON(c, "Only admin users can access this endpoint", http.StatusForbidden, nil, nil)
            return
        }

        productIDStr := c.Param("product_id")
        if productIDStr == "" {
            response.JSON(c, "Product ID cannot be empty", http.StatusBadRequest, nil, nil)
            return
        }

        productID64, err := strconv.ParseUint(productIDStr, 10, 32)
        if err != nil {
            response.JSON(c, "Invalid product ID", http.StatusBadRequest, nil, err)
            return
        }
        
        productID := uint(productID64)
        var product models.Product
        if err := c.ShouldBindJSON(&product); err != nil {
            response.JSON(c, "Invalid JSON format", http.StatusBadRequest, nil, err)
            return
        }

        product.ID = productID
        if err := s.ProductRepo.UpdateProduct(&product); err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                response.JSON(c, "Product not found", http.StatusNotFound, nil, nil)
                return
            }
            response.JSON(c, "Failed to update product", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Product updated successfully", http.StatusOK, nil, nil)
    }
}

// handleDeleteProduct deletes a product by ID
// @Summary Delete a product by ID
// @Description Delete a product from the inventory (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 403 {object} response.ErrorResponse "Only admin users can access this endpoint"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /products/{product_id} [delete]
func (s *Server) handleDeleteProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        if userRole != "Admin" {
            response.JSON(c, "Only admin users can access this endpoint", http.StatusForbidden, nil, nil)
            return
        }

        productIDStr := c.Param("product_id")
        if productIDStr == "" {
            response.JSON(c, "Product ID cannot be empty", http.StatusBadRequest, nil, nil)
            return
        }

        productID64, err := strconv.ParseUint(productIDStr, 10, 32)
        if err != nil {
            response.JSON(c, "Invalid product ID", http.StatusBadRequest, nil, err)
            return
        }
        
        productID := uint(productID64)
        if err := s.ProductRepo.DeleteProduct(productID); err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                response.JSON(c, "Product not found", http.StatusNotFound, nil, nil)
                return
            }
            response.JSON(c, "Failed to delete product", http.StatusInternalServerError, nil, err)
            return
        }

        c.Status(http.StatusNoContent)
    }
}
