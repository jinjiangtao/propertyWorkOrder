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

type AssignWorkerRequest struct {
	RepairID  int64 `json:"repair_id" binding:"required"`
	WorkerID  int64 `json:"worker_id" binding:"required"`
}

type WorkerActionRequest struct {
	RepairID int64 `json:"repair_id" binding:"required"`
	WorkerID int64 `json:"worker_id" binding:"required"`
}

type RepairResultRequest struct {
	RepairID     int64  `json:"repair_id" binding:"required"`
	WorkerID     int64  `json:"worker_id" binding:"required"`
	RepairResult string `json:"repair_result" binding:"required"`
	RepairImgs   string `json:"repair_imgs"`
}

func UpdateRepairStatus(c *gin.Context) {
	var req UpdateRepairStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID和状态"})
		return
	}

	validStatuses := map[string]bool{
		"未处理": true,
		"已派单": true,
		"处理中": true,
		"已完成": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态，只能是：未处理、已派单、处理中、已完成"})
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

func AssignWorker(c *gin.Context) {
	var req AssignWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID和工人ID"})
		return
	}

	var workerName string
	err := DB.QueryRow("SELECT name FROM workers WHERE id = ? AND status = 1", req.WorkerID).Scan(&workerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工人不存在或已离职"})
		return
	}

	result, err := DB.Exec(
		"UPDATE repair_requests SET worker_id = ?, worker_name = ?, status = '已派单', updated_at = ? WHERE id = ?",
		req.WorkerID,
		workerName,
		time.Now(),
		req.RepairID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "派单失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "派单成功",
	})
}

func WorkerAcceptOrder(c *gin.Context) {
	var req WorkerActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID和工人ID"})
		return
	}

	var currentStatus string
	var currentWorkerID int64
	err := DB.QueryRow(
		"SELECT status, worker_id FROM repair_requests WHERE id = ?",
		req.RepairID,
	).Scan(&currentStatus, &currentWorkerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	if currentStatus != "已派单" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能接已派单的工单"})
		return
	}

	if currentWorkerID != req.WorkerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能接分配给自己的工单"})
		return
	}

	result, err := DB.Exec(
		"UPDATE repair_requests SET status = '处理中', updated_at = ? WHERE id = ?",
		time.Now(),
		req.RepairID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "接单失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "接单成功",
	})
}

func WorkerRejectOrder(c *gin.Context) {
	var req WorkerActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID和工人ID"})
		return
	}

	var currentStatus string
	var currentWorkerID int64
	err := DB.QueryRow(
		"SELECT status, worker_id FROM repair_requests WHERE id = ?",
		req.RepairID,
	).Scan(&currentStatus, &currentWorkerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	if currentStatus != "已派单" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能拒绝已派单的工单"})
		return
	}

	if currentWorkerID != req.WorkerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能拒绝分配给自己的工单"})
		return
	}

	result, err := DB.Exec(
		"UPDATE repair_requests SET status = '未处理', worker_id = 0, worker_name = '', updated_at = ? WHERE id = ?",
		time.Now(),
		req.RepairID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "拒单失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "拒单成功，工单已退回未处理状态",
	})
}

func SubmitRepairResult(c *gin.Context) {
	var req RepairResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供报修ID、工人ID和维修结果"})
		return
	}

	var currentStatus string
	var currentWorkerID int64
	err := DB.QueryRow(
		"SELECT status, worker_id FROM repair_requests WHERE id = ?",
		req.RepairID,
	).Scan(&currentStatus, &currentWorkerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	if currentStatus != "处理中" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能提交处理中工单的维修结果"})
		return
	}

	if currentWorkerID != req.WorkerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能提交分配给自己的工单"})
		return
	}

	result, err := DB.Exec(
		"UPDATE repair_requests SET status = '已完成', repair_result = ?, repair_imgs = ?, updated_at = ? WHERE id = ?",
		req.RepairResult,
		req.RepairImgs,
		time.Now(),
		req.RepairID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交维修结果失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "报修记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "维修结果提交成功",
	})
}

func GetWorkerRepairs(c *gin.Context) {
	workerIDStr := c.Query("worker_id")
	if workerIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供工人ID"})
		return
	}

	workerID, err := strconv.ParseInt(workerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工人ID"})
		return
	}

	rows, err := DB.Query(
		`SELECT id, user_id, username, repair_type, description, image_url, status, created_at, updated_at, worker_id, worker_name, repair_result, repair_imgs
		FROM repair_requests
		WHERE worker_id = ?
		ORDER BY created_at DESC`,
		workerID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询工单失败"})
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
			&repair.WorkerID,
			&repair.WorkerName,
			&repair.RepairResult,
			&repair.RepairImgs,
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

func GetWorkerStats(c *gin.Context) {
	rows, err := DB.Query(`
		SELECT w.id, w.name, w.work_no, w.skill_type, 
			COUNT(r.id) as total_count,
			SUM(CASE WHEN r.status = '已完成' THEN 1 ELSE 0 END) as completed_count
		FROM workers w
		LEFT JOIN repair_requests r ON w.id = r.worker_id
		GROUP BY w.id, w.name, w.work_no, w.skill_type
		ORDER BY completed_count DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询统计数据失败"})
		return
	}
	defer rows.Close()

	type WorkerStats struct {
		ID             int64  `json:"id"`
		Name           string `json:"name"`
		WorkNo         string `json:"work_no"`
		SkillType      string `json:"skill_type"`
		TotalCount     int    `json:"total_count"`
		CompletedCount int    `json:"completed_count"`
	}

	var stats []WorkerStats
	for rows.Next() {
		var stat WorkerStats
		if err := rows.Scan(
			&stat.ID,
			&stat.Name,
			&stat.WorkNo,
			&stat.SkillType,
			&stat.TotalCount,
			&stat.CompletedCount,
		); err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
