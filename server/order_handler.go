package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/techagentng/ecommerce-api/models"
	"github.com/techagentng/ecommerce-api/server/response"
)

// handlePlaceOrder handles placing an order
func (s *Server) handlePlaceOrder() gin.HandlerFunc {
    return func(c *gin.Context) {
        var orderRequest struct {
            Items []struct {
                ProductID uint `json:"product_id" binding:"required"`
                Quantity  int  `json:"quantity" binding:"required"`
            } `json:"items" binding:"required"`
        }

        // Bind the JSON request to the orderRequest struct
        if err := c.ShouldBindJSON(&orderRequest); err != nil {
            response.JSON(c, "Invalid order request", http.StatusBadRequest, nil, err)
            return
        }

        // Retrieve user ID from context
        userID := c.GetUint("userID")
        log.Printf("UserID00000000000000000000000: %d", userID)
        if userID == 0 {
            response.JSON(c, "User not authenticated", http.StatusUnauthorized, nil, nil)
            return
        }

        // Check if the user exists
        var user models.User
        if err := s.OrderRepo.FindUserByID(userID, &user); err != nil || user.ID == 0 {
            response.JSON(c, "User not found", http.StatusBadRequest, nil, nil)
            return
        }

        var totalOrderPrice float64
        var orderItems []models.OrderItem

        // Iterate through each item in the order request
        for _, item := range orderRequest.Items {
            product, err := s.ProductRepo.FindProductByID(item.ProductID)
            if err != nil {
                response.JSON(c, "Error retrieving product", http.StatusInternalServerError, nil, err)
                return
            }

            if product == nil {
                response.JSON(c, "Product not found", http.StatusBadRequest, nil, nil)
                return
            }

            // Calculate the total price for this item
            itemTotal := float64(item.Quantity) * product.Price
            totalOrderPrice += itemTotal

            // Add the item to the orderItems slice
            orderItems = append(orderItems, models.OrderItem{
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
                UnitPrice: product.Price,
                TotalPrice: itemTotal,
            })
        }

        // Create the Order
        order := models.Order{
            UserID:     userID, 
            TotalPrice: totalOrderPrice,
            Status:     "Pending",
            Items:      orderItems,
        }

        // Save the order
        createdOrder, err := s.OrderRepo.CreateOrder(&order)
        if err != nil {
            response.JSON(c, "Failed to place order", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Order placed successfully", http.StatusOK, createdOrder, nil)
    }
}

// handleUpdateOrderStatus updates the status of an order (admin privilege)
func (s *Server) handleUpdateOrderStatus() gin.HandlerFunc {
    return func(c *gin.Context) {
        orderID := c.Param("id")
        parsedOrderID, err := uuid.Parse(orderID)
        if err != nil {
            response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
            return
        }

        // Assume you have a structure to capture the incoming status update
        var statusUpdate struct {
            Status string `json:"status"`
        }

        if err := c.ShouldBindJSON(&statusUpdate); err != nil {
            response.JSON(c, "Invalid request payload", http.StatusBadRequest, nil, err)
            return
        }

        // Call UpdateOrderStatus and capture both return values
        updatedOrder, err := s.OrderService.UpdateOrderStatus(parsedOrderID, statusUpdate.Status)
        if err != nil {
            response.JSON(c, "Failed to update order status", http.StatusInternalServerError, nil, err)
            return
        }

        // Return the updated order in the response
        response.JSON(c, "Order status updated successfully", http.StatusOK, updatedOrder, nil)
    }
}


