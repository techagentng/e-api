package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	apiError "github.com/techagentng/ecommerce-api/errors"
	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/gorm"
)

// OrderService interface
type OrderService interface {
	PlaceOrder(order *models.Order) (*models.Order, error)
	ListUserOrders(userID uuid.UUID) ([]models.Order, error)
	CancelOrder(orderID uuid.UUID) (*models.Order, error)
	UpdateOrderStatus(orderID uuid.UUID, status string) (*models.Order, error)
}

// orderService struct
type orderService struct {
	Config    *config.Config
	orderRepo db.OrderRepository
}

// NewOrderService constructor function
func NewOrderService(orderRepo db.OrderRepository, conf *config.Config) OrderService {
	return &orderService{
		Config:    conf,
		orderRepo: orderRepo,
	}
}

// PlaceOrder allows a user to place a new order
func (o *orderService) PlaceOrder(order *models.Order) (*models.Order, error) {
	if err,_ := o.orderRepo.CreateOrder(order); err != nil {
		log.Printf("Error placing order: %v", err)
		return nil, apiError.New("unable to place order", http.StatusInternalServerError)
	}
	return order, nil
}

// ListUserOrders retrieves all orders for a specific user
func (o *orderService) ListUserOrders(userID uuid.UUID) ([]models.Order, error) {
    orders, err := o.orderRepo.FindOrdersByUserID(userID) 
    if err != nil {
        log.Printf("Error fetching orders for user %s: %v", userID, err)
        return nil, apiError.New("unable to fetch user orders", http.StatusInternalServerError)
    }

    var ordersResponse []models.Order
    for _, order := range orders {
        ordersResponse = append(ordersResponse, *order)
    }

    return ordersResponse, nil
}

// UpdateOrderStatus allows an admin to update the status of an order
func (o *orderService) UpdateOrderStatus(orderID uuid.UUID, status string) (*models.Order, error) {
	order, err := o.orderRepo.GetOrderByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiError.New("order not found", http.StatusNotFound)
		}
		log.Printf("Error fetching order: %v", err)
		return nil, apiError.New("unable to fetch order", http.StatusInternalServerError)
	}

	order.Status = status
    if err := o.orderRepo.UpdateOrderStatus(order.ID, status); err != nil {
        log.Printf("Error updating order status: %v", err)
        return nil, apiError.New("unable to update order status", http.StatusInternalServerError)
    }
	return order, nil
}

// Service function to cancel an order
func (o *orderService) CancelOrder(orderID uuid.UUID) (*models.Order, error) {
	order, err := o.orderRepo.FindOrderByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiError.New("order not found", http.StatusNotFound)
		}
		log.Printf("Error retrieving order for cancellation: %v", err)
		return nil, apiError.New("unable to retrieve order for cancellation", http.StatusInternalServerError)
	}

	if order.Status != "Pending" {
		return nil, apiError.New("only pending orders can be canceled", http.StatusBadRequest)
	}

	// Update the order status to 'Canceled'
	order.Status = "Canceled"
	if err := o.orderRepo.UpdateOrder(order); err != nil {
		log.Printf(
			"Error updating order status to 'Canceled' for order ID %v. Details: %v",
			order.ID,
			err,
		)
		return nil, apiError.New("unable to cancel order due to a database error", http.StatusInternalServerError)
	}

	log.Printf("Order with ID %v successfully canceled", orderID)
	return order, nil
}

