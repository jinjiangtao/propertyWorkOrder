package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateRepairRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Username    string `json:"username" binding:"required"`
	RepairType  string `json:"repair_type" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url"`
}

func CreateRepair(c *gin.Context) {
	var req CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供完整的报修信息"})
		return
	}

	result, err := DB.Exec(
		`INSERT INTO repair_requests (user_id, username, repair_type, description, image_url, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, '未处理', ?, ?)`,
		req.UserID,
		req.Username,
		req.RepairType,
		req.Description,
		req.ImageURL,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建报修单失败"})
		return
	}

	requestID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "报修单创建成功",
		"request_id": requestID,
	})
}

func GetUserRepairs(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户ID"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	rows, err := DB.Query(
		`SELECT id, user_id, username, repair_type, description, image_url, status, created_at, updated_at
		FROM repair_requests
		WHERE user_id = ?
		ORDER BY created_at DESC`,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询报修记录失败"})
		return
	}
	defer rows.Close()

	var repairs []RepairRequest
	for rows.Next() {
		var repair RepairRequest
		if err := rows.Scan(
			&repair.ID,
			&repair.UserID,
			&repair.Username,
			&repair.RepairType,
			&repair.Description,
			&repair.ImageURL,
			&repair.Status,
			&repair.CreatedAt,
			&repair.UpdatedAt,
		); err != nil {
			continue
		}
		repairs = append(repairs, repair)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    repairs,
	})
}

func GetAllRepairs(c *gin.Context) {
	rows, err := DB.Query(
		`SELECT id, user_id, username, repair_type, description, image_url, status, created_at, updated_at
		FROM repair_requests
		ORDER BY created_at DESC`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询报修记录失败"})
		return
	}
	defer rows.Close()

	var repairs []RepairRequest
	for rows.Next() {
		var repair RepairRequest
		if err := rows.Scan(
			&repair.ID,
			&repair.UserID,
			&repair.Username,
			&repair.RepairType,
			&repair.Description,
			&repair.ImageURL,
			&repair.Status,
			&repair.CreatedAt,
			&repair.UpdatedAt,
		); err != nil {
			continue
		}
		repairs = append(repairs, repair)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    repairs,
	})
}

type UpdateRepairStatusRequest struct {
	RepairID int64  `json:"repair_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

func UpdateRepairStatus(c *gin.Context) {
	var req UpdateRepairStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID和状态"})
		return
	}

	validStatuses := map[string]bool{
		"未处理": true,
		"处理中": true,
		"已完成": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态，只能是：未处理、处理中、已完成"})
		return
	}

	result, err := DB.Exec(
		"UPDATE repair_requests SET status = ?, updated_at = ? WHERE id = ?",
		req.Status,
		time.Now(),
		req.RepairID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "状态更新成功",
	})
}
