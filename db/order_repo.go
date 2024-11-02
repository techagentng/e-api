package db

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	apiError "github.com/techagentng/ecommerce-api/errors"
	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(orderRequest *models.Order) ([]*models.Order, error)
	FindOrderByID(id uuid.UUID) (*models.Order, error)
	FindOrdersByUserID(userID uint) ([]*models.Order, error)
	UpdateOrderStatus(id uint, status string) error
	CancelOrder(id uint) error
	UpdateOrder(order *models.Order) error
	GetOrderByID(orderID uuid.UUID) (*models.Order, error)
	GetOrdersByUserID(userID uint) ([]*models.Order, error)
	GetUserIDFromUUID(userUUID uuid.UUID) (uint, error)
	FindUserByID(userID uint, user *models.User) error
	LoadOrderDetails(orderID uint) (*models.Order, error)
}

type orderRepo struct {
	DB *gorm.DB
}

func NewOrderRepo(db *GormDB) OrderRepository {
	return &orderRepo{db.DB}
}

func (o *orderRepo) CreateOrder(orderRequest *models.Order) ([]*models.Order, error) {
    var createdOrders []*models.Order
    for _, item := range orderRequest.Items {
        var product models.Product
        
        if err := o.DB.First(&product, item.ProductID).Error; err != nil {
            return nil, err 
        }

        order := &models.Order{
            UserID:     orderRequest.UserID, 
            ProductID:  item.ProductID,
            Quantity:   item.Quantity,
            TotalPrice: float64(item.Quantity) * product.Price, 
            Status:     "Pending", 
        }

        // Save the order to the database
        if err := o.DB.Create(order).Error; err != nil {
            return nil, err 
        }

        // Populate related user data
        var user models.User
        if err := o.DB.First(&user, order.UserID).Error; err == nil {
            order.User = user 
        } else {
            order.User = models.User{} 
        }

        order.Product = product
        createdOrders = append(createdOrders, order)
    }

    return createdOrders, nil
}

func (o *orderRepo) FindOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := o.DB.First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (o *orderRepo) UpdateOrderStatus(id uint, status string) error {
	return o.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (o *orderRepo) CancelOrder(id uint) error {
	return o.DB.Model(&models.Order{}).Where("id = ? AND status = ?", id, "Pending").Update("status", "Cancelled").Error
}

func (o *orderRepo) UpdateOrder(order *models.Order) error {
	// Perform the update operation on the order
	if err := o.DB.Save(order).Error; err != nil {
		log.Printf("Error updating order with ID %v: %v", order.ID, err)
		return err
	}

	log.Printf("Order with ID %v successfully updated", order.ID)
	return nil
}

func (o *orderRepo) GetOrderByID(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := o.DB.First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err 
	}
	return &order, nil 
}

func (o *orderRepo) GetOrdersByUserID(userID uint) ([]*models.Order, error) {
	var orders []*models.Order

	// Perform the query to fetch orders by userID
	if err := o.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err // Return the error for further handling
	}

	return orders, nil
}

func (o *orderRepo) GetOrdersForUser(userID uint) ([]*models.Order, error) {
	orders, err := o.GetOrdersByUserID(userID)
	if err != nil {
		log.Printf("Error retrieving orders for user ID %v: %v", userID, err)
		return nil, apiError.New("unable to retrieve orders for the user", http.StatusInternalServerError)
	}

	if len(orders) == 0 {
		log.Printf("No orders found for user ID %v", userID)
		return nil, apiError.New("no orders found for the specified user", http.StatusNotFound)
	}
	return orders, nil
}

func (o *orderRepo)  GetUserIDFromUUID(userUUID uuid.UUID) (uint, error) {
	var user models.User

	// Query the database for the user with the given UUID
	if err := o.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, err // Return the actual error if it's something else
	}

	// Return the user ID
	return user.ID, nil
}

func (o *orderRepo) FindUserByID(userID uint, user *models.User) error {
    if err := o.DB.First(user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return gorm.ErrRecordNotFound // User not found
        }
        return err // Return any other errors
    }
    return nil // User found successfully
}

func (o *orderRepo) FindOrdersByUserID(userID uint) ([]*models.Order, error) {
	var orders []*models.Order
	if err := o.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderRepo) LoadOrderDetails(orderID uint) (*models.Order, error) {
    var order models.Order
    if err := o.DB.Preload("User").Preload("Product").First(&order, "id = ?", orderID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil 
        }
        return nil, err 
    }

    return &order, nil 
}