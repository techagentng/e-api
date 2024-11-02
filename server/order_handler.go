package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

        if err := c.ShouldBindJSON(&orderRequest); err != nil {
            response.JSON(c, "Invalid order request", http.StatusBadRequest, nil, err)
            return
        }

        userID := c.GetUint("userID")
        if userID == 0 {
            response.JSON(c, "User not authenticated", http.StatusUnauthorized, nil, nil)
            return
        }

        var user models.User
        if err := s.OrderRepo.FindUserByID(userID, &user); err != nil || user.ID == 0 {
            response.JSON(c, "User not found", http.StatusBadRequest, nil, nil)
            return
        }

        var totalOrderPrice float64
        var orderItems []models.OrderItem
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

            itemTotal := float64(item.Quantity) * product.Price
            totalOrderPrice += itemTotal
            orderItems = append(orderItems, models.OrderItem{
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
                UnitPrice: product.Price,
                TotalPrice: itemTotal,
            })
        }

        order := models.Order{
            UserID:     userID,
            TotalPrice: totalOrderPrice,
            Status:     "Pending",
            Items:      orderItems,
        }

        createdOrders, err := s.OrderRepo.CreateOrder(&order)
        if err != nil {
            response.JSON(c, "Failed to place order", http.StatusInternalServerError, nil, err)
            return
        }

        // Load order details for each created order
        var ordersWithDetails []*models.Order 
        for _, createdOrder := range createdOrders {
            orderWithDetails, err := s.OrderRepo.LoadOrderDetails(createdOrder.ID)
            if err != nil {
                response.JSON(c, "Failed to load order details", http.StatusInternalServerError, nil, err)
                return
            }
            if orderWithDetails != nil { // Check for nil before appending
                ordersWithDetails = append(ordersWithDetails, orderWithDetails)
            }
        }

        // Create response DTO
        responseDTO := models.PlaceOrderResponse{
            OrderID:    ordersWithDetails[0].ID, // Assuming you want the ID of the first order
            UserID:     userID,
            TotalPrice: totalOrderPrice,
            Status:     ordersWithDetails[0].Status,
            CreatedAt:  ordersWithDetails[0].CreatedAt.Format(time.RFC3339),
        }

        response.JSON(c, "Order placed successfully", http.StatusOK, responseDTO, nil)
    }
}

func (s *Server) handleListUserOrders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve userID from the context, assuming it's set by middleware
        userID := c.GetUint("userID")
        log.Printf("UserID: %d", userID)
        
        // Check if the userID is valid
        if userID == 0 {
            response.JSON(c, "User not authenticated", http.StatusUnauthorized, nil, nil)
            return
        }

        // Find user by userID to validate existence
        var user models.User
        if err := s.OrderRepo.FindUserByID(userID, &user); err != nil || user.ID == 0 {
            response.JSON(c, "User not found", http.StatusBadRequest, nil, nil)
            return
        }

        // Retrieve orders for the given userID
        orders, err := s.OrderRepo.FindOrdersByUserID(userID)
        if err != nil {
            log.Printf("Error fetching orders: %v", err)
            response.JSON(c, "Failed to fetch orders", http.StatusInternalServerError, nil, err)
            return
        }

        if len(orders) == 0 {
            response.JSON(c, "No orders found", http.StatusOK, nil, nil)
            return
        }

        response.JSON(c, "Orders retrieved successfully", http.StatusOK, orders, nil)
    }
}

func (s *Server) handleCancelOrder() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        userID, _ := c.Get("userID")
        orderIDStr := c.Param("order_id")
        fmt.Println("Order ID parameter:", orderIDStr) 
        if orderIDStr == "" {
            response.JSON(c, "Order ID cannot be empty", http.StatusBadRequest, nil, nil)
            return
        }
        orderID64, err := strconv.ParseUint(orderIDStr, 10, 32) 
        if err != nil {
            response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
            return
        }
        
        orderID := uint(orderID64) 
        order, err := s.OrderRepo.LoadOrderDetails(orderID)
        if err != nil {
            response.JSON(c, "Failed to load order details", http.StatusInternalServerError, nil, err)
            return
        }
        if order == nil {
            response.JSON(c, "Order not found", http.StatusNotFound, nil, nil)
            return
        }

        if order.Status != "Pending" {
            response.JSON(c, "Only pending orders can be canceled", http.StatusBadRequest, nil, nil)
            return
        }

        if userRole != "Admin" && order.UserID != userID {
            response.JSON(c, "Access denied: You cannot cancel this order", http.StatusForbidden, nil, nil)
            return
        }

        // Update the order status to "Canceled"
        err = s.OrderRepo.CancelOrder(orderID)
        if err != nil {
            response.JSON(c, "Failed to cancel order", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Order canceled successfully", http.StatusOK, nil, nil)
    }
}

type UpdateStatusRequest struct {
    Status string `json:"status"`
}

func (s *Server) handleUpdateOrderStatus() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("user_role")
        if userRole != "Admin" {
            response.JSON(c, "Only admin users can update the order status", http.StatusForbidden, nil, nil)
            return
        }

        orderIDStr := c.Param("order_id")
        if orderIDStr == "" {
            response.JSON(c, "Order ID cannot be empty", http.StatusBadRequest, nil, nil)
            return
        }

        orderID64, err := strconv.ParseUint(orderIDStr, 10, 32)
        if err != nil {
            response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
            return
        }
        
        orderID := uint(orderID64)
        order, err := s.OrderRepo.LoadOrderDetails(orderID)
        if err != nil {
            response.JSON(c, "Failed to load order details", http.StatusInternalServerError, nil, err)
            return
        }
        if order == nil {
            response.JSON(c, "Order not found", http.StatusNotFound, nil, nil)
            return
        }

        var req UpdateStatusRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            response.JSON(c, "Invalid JSON format", http.StatusBadRequest, nil, err)
            return
        }
        newStatus := req.Status

        allowedStatuses := []string{"Pending", "Canceled", "Completed", "Shipped"}
        isValidStatus := false
        for _, status := range allowedStatuses {
            if newStatus == status {
                isValidStatus = true
                break
            }
        }
        if !isValidStatus {
            response.JSON(c, "Invalid status value", http.StatusBadRequest, nil, nil)
            return
        }

        err = s.OrderRepo.UpdateOrderStatus(orderID, newStatus)
        if err != nil {
            response.JSON(c, "Failed to update order status", http.StatusInternalServerError, nil, err)
            return
        }

        response.JSON(c, "Order status updated successfully", http.StatusOK, nil, nil)
    }
}


