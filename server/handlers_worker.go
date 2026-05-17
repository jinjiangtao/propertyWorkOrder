package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type WorkerLoginRequest struct {
	WorkNo   string `json:"work_no" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WorkerLoginResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	WorkerID int64  `json:"worker_id,omitempty"`
	Name     string `json:"name,omitempty"`
	SkillType string `json:"skill_type,omitempty"`
}

type CreateWorkerRequest struct {
	WorkNo    string `json:"work_no" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required"`
	SkillType string `json:"skill_type" binding:"required"`
}

type UpdateWorkerRequest struct {
	ID        int64  `json:"id" binding:"required"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	SkillType string `json:"skill_type"`
}

func WorkerLogin(c *gin.Context) {
	var req WorkerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供工号和密码"})
		return
	}

	var worker Worker
	err := DB.QueryRow(
		"SELECT id, work_no, name, phone, password, skill_type, status FROM workers WHERE work_no = ?",
		req.WorkNo,
	).Scan(&worker.ID, &worker.WorkNo, &worker.Name, &worker.Phone, &worker.Password, &worker.SkillType, &worker.Status)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "工号或密码错误"})
		return
	}

	if worker.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "该工人已离职"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(worker.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "工号或密码错误"})
		return
	}

	response := WorkerLoginResponse{
		Success:  true,
		Message:  "登录成功",
		WorkerID: worker.ID,
		Name:     worker.Name,
		SkillType: worker.SkillType,
	}

	c.JSON(http.StatusOK, response)
}

func GetWorkers(c *gin.Context) {
	rows, err := DB.Query(
		"SELECT id, work_no, name, phone, skill_type, status, created_at FROM workers ORDER BY created_at DESC",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询工人列表失败"})
		return
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var worker Worker
		if err := rows.Scan(
			&worker.ID,
			&worker.WorkNo,
			&worker.Name,
			&worker.Phone,
			&worker.SkillType,
			&worker.Status,
			&worker.CreatedAt,
		); err != nil {
			continue
		}
		workers = append(workers, worker)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    workers,
	})
}

func CreateWorker(c *gin.Context) {
	var req CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供完整的工人信息"})
		return
	}

	validSkillTypes := map[string]bool{
		"水电":     true,
		"木工":     true,
		"保洁":     true,
		"综合维修": true,
	}
	if !validSkillTypes[req.SkillType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工种类型，只能是：水电、木工、保洁、综合维修"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	result, err := DB.Exec(
		"INSERT INTO workers (work_no, name, phone, password, skill_type, status, created_at) VALUES (?, ?, ?, ?, ?, 1, ?)",
		req.WorkNo,
		req.Name,
		req.Phone,
		string(hashedPassword),
		req.SkillType,
		time.Now(),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工号或手机号已存在"})
		return
	}

	workerID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "工人创建成功",
		"worker_id": workerID,
	})
}

func UpdateWorker(c *gin.Context) {
	var req UpdateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供工人ID"})
		return
	}

	if req.SkillType != "" {
		validSkillTypes := map[string]bool{
			"水电":     true,
			"木工":     true,
			"保洁":     true,
			"综合维修": true,
		}
		if !validSkillTypes[req.SkillType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工种类型"})
			return
		}
	}

	var updates []interface{}
	var setClauses []string

	if req.Name != "" {
		setClauses = append(setClauses, "name = ?")
		updates = append(updates, req.Name)
	}
	if req.Phone != "" {
		setClauses = append(setClauses, "phone = ?")
		updates = append(updates, req.Phone)
	}
	if req.SkillType != "" {
		setClauses = append(setClauses, "skill_type = ?")
		updates = append(updates, req.SkillType)
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有需要更新的字段"})
		return
	}

	updates = append(updates, req.ID)
	query := "UPDATE workers SET " + join(setClauses, ", ") + " WHERE id = ?"

	result, err := DB.Exec(query, updates...)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号已被使用"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "工人不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "工人信息更新成功",
	})
}

func ToggleWorkerStatus(c *gin.Context) {
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

	var currentStatus int
	err = DB.QueryRow("SELECT status FROM workers WHERE id = ?", workerID).Scan(&currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工人不存在"})
		return
	}

	newStatus := 1
	if currentStatus == 1 {
		newStatus = 0
	}

	result, err := DB.Exec("UPDATE workers SET status = ? WHERE id = ?", newStatus, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "工人不存在"})
		return
	}

	statusText := "启用"
	if newStatus == 0 {
		statusText = "禁用"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "工人已" + statusText,
	})
}

func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}