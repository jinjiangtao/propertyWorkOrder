package service

import (
	"backend/entity"
	"backend/model"
)

func CreateWorkOrder(order *entity.WorkOrder) error {
	return model.CreateWorkOrder(order)
}

func GetWorkOrdersByUserID(userID int) ([]entity.WorkOrder, error) {
	return model.GetWorkOrdersByUserID(userID)
}

func GetAllWorkOrders() ([]entity.WorkOrder, error) {
	return model.GetAllWorkOrders()
}

func UpdateWorkOrderStatus(orderID int, status int) error {
	return model.UpdateWorkOrderStatus(orderID, status)
}

func GetWorkOrderByID(orderID int) (*entity.WorkOrder, error) {
	return model.GetWorkOrderByID(orderID)
}