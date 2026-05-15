package model

import (
	"backend/entity"
	"errors"
	"sync"
)

var users []entity.User
var workOrders []entity.WorkOrder
var userIDCounter = 1
var orderIDCounter = 1
var mu sync.Mutex

func InitDB() error {
	mu.Lock()
	defer mu.Unlock()

	adminExists := false
	for _, u := range users {
		if u.Username == "admin" {
			adminExists = true
			break
		}
	}

	if !adminExists {
		users = append(users, entity.User{
			ID:       userIDCounter,
			Username: "admin",
			Password: "admin123",
			Role:     entity.RoleAdmin,
		})
		userIDCounter++
	}

	return nil
}

func GetUserByUsername(username string) (*entity.User, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func CreateUser(user *entity.User) error {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	user.ID = userIDCounter
	userIDCounter++
	users = append(users, *user)
	return nil
}

func CheckUsernameExists(username string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func CreateWorkOrder(order *entity.WorkOrder) error {
	mu.Lock()
	defer mu.Unlock()

	order.ID = orderIDCounter
	orderIDCounter++
	workOrders = append(workOrders, *order)
	return nil
}

func GetWorkOrdersByUserID(userID int) ([]entity.WorkOrder, error) {
	mu.Lock()
	defer mu.Unlock()

	var result []entity.WorkOrder
	for _, o := range workOrders {
		if o.UserId == userID {
			result = append(result, o)
		}
	}
	return result, nil
}

func GetAllWorkOrders() ([]entity.WorkOrder, error) {
	mu.Lock()
	defer mu.Unlock()

	return workOrders, nil
}

func UpdateWorkOrderStatus(orderID int, status int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, o := range workOrders {
		if o.ID == orderID {
			workOrders[i].Status = status
			return nil
		}
	}
	return errors.New("order not found")
}

func GetWorkOrderByID(orderID int) (*entity.WorkOrder, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, o := range workOrders {
		if o.ID == orderID {
			return &o, nil
		}
	}
	return nil, errors.New("order not found")
}