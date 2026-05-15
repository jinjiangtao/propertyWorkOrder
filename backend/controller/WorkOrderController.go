package controller

import (
	"backend/entity"
	"backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateWorkOrderRequest struct {
	UserId      int    `json:"user_id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Images      string `json:"images"`
}

func CreateWorkOrder(c *gin.Context) {
	var req CreateWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	order := &entity.WorkOrder{
		UserId:      req.UserId,
		Type:        req.Type,
		Description: req.Description,
		Images:      req.Images,
		Status:      entity.StatusPending,
	}

	if err := service.CreateWorkOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func GetWorkOrders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		orders, err := service.GetAllWorkOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	orders, err := service.GetWorkOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

type UpdateStatusRequest struct {
	Status int `json:"status"`
}

func UpdateWorkOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status != entity.StatusPending && req.Status != entity.StatusProcessing && req.Status != entity.StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
		return
	}

	if err := service.UpdateWorkOrderStatus(orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}