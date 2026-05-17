package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) {
	os.Remove("./test.db")
	var err error
	DB, err = sql.Open("sqlite", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	if err = createTables(); err != nil {
		t.Fatal(err)
	}
	if err = createAdminUser(); err != nil {
		t.Fatal(err)
	}
}

func cleanupTestDB() {
	DB.Close()
	os.Remove("./test.db")
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/register", Register)
	r.POST("/api/login", Login)
	r.POST("/api/repair/create", CreateRepair)
	r.GET("/api/repair/user", GetUserRepairs)
	r.GET("/api/repair/all", GetAllRepairs)
	r.PUT("/api/repair/status", UpdateRepairStatus)
	r.POST("/api/repair/assign", AssignWorker)
	r.POST("/api/repair/accept", WorkerAcceptOrder)
	r.POST("/api/repair/reject", WorkerRejectOrder)
	r.POST("/api/repair/result", SubmitRepairResult)
	
	r.POST("/api/worker/login", WorkerLogin)
	r.GET("/api/worker/list", GetWorkers)
	r.POST("/api/worker/create", CreateWorker)
	r.PUT("/api/worker/update", UpdateWorker)
	r.PUT("/api/worker/status", ToggleWorkerStatus)
	return r
}

func TestRegister(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid registration",
			body:       map[string]string{"username": "testuser", "password": "testpass"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing username",
			body:       map[string]string{"password": "testpass"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       map[string]string{"username": "testuser"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate username",
			body:       map[string]string{"username": "testuser", "password": "testpass"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid login",
			body:       map[string]string{"username": "testuser", "password": "testpass"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong password",
			body:       map[string]string{"username": "testuser", "password": "wrongpass"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "non-existent user",
			body:       map[string]string{"username": "nonexistent", "password": "testpass"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestCreateRepair(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid repair request",
			body: map[string]interface{}{
				"user_id":     loginResp.UserID,
				"username":    "testuser",
				"repair_type": "水电维修",
				"description": "水龙头漏水",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing repair type",
			body: map[string]interface{}{
				"user_id":     loginResp.UserID,
				"username":    "testuser",
				"description": "水龙头漏水",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestGetUserRepairs(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	repairBody := map[string]interface{}{
		"user_id":     loginResp.UserID,
		"username":    "testuser",
		"repair_type": "水电维修",
		"description": "水龙头漏水",
	}
	jsonBody, _ = json.Marshal(repairBody)
	req, _ = http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	url := fmt.Sprintf("/api/repair/user?user_id=%d", loginResp.UserID)
	req, _ = http.NewRequest("GET", url, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestGetAllRepairs(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/repair/all", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestUpdateRepairStatus(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	repairBody := map[string]interface{}{
		"user_id":     loginResp.UserID,
		"username":    "testuser",
		"repair_type": "水电维修",
		"description": "水龙头漏水",
	}
	jsonBody, _ = json.Marshal(repairBody)
	req, _ = http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid status update to 处理中",
			body: map[string]interface{}{
				"repair_id": 1,
				"status":    "处理中",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "valid status update to 已完成",
			body: map[string]interface{}{
				"repair_id": 1,
				"status":    "已完成",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid status",
			body: map[string]interface{}{
				"repair_id": 1,
				"status":    "无效状态",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("PUT", "/api/repair/status", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestCreateWorker(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid worker creation",
			body: map[string]string{"work_no": "W001", "name": "张三", "phone": "13800138001", "password": "123456", "skill_type": "水电"},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing work_no",
			body: map[string]string{"name": "李四", "phone": "13800138002", "password": "123456", "skill_type": "木工"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate work_no",
			body: map[string]string{"work_no": "W001", "name": "王五", "phone": "13800138003", "password": "123456", "skill_type": "保洁"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/worker/create", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestWorkerLogin(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	body := map[string]string{"work_no": "W001", "name": "张三", "phone": "13800138001", "password": "123456", "skill_type": "水电"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/worker/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid worker login",
			body: map[string]string{"work_no": "W001", "password": "123456"},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: map[string]string{"work_no": "W001", "password": "wrongpass"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent worker",
			body: map[string]string{"work_no": "W999", "password": "123456"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/worker/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestAssignWorker(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	registerBody := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	repairBody := map[string]interface{}{
		"user_id":     loginResp.UserID,
		"username":    "testuser",
		"repair_type": "水电维修",
		"description": "水龙头漏水",
	}
	jsonBody, _ = json.Marshal(repairBody)
	req, _ = http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	workerBody := map[string]string{"work_no": "W001", "name": "张三", "phone": "13800138001", "password": "123456", "skill_type": "水电"}
	jsonBody, _ = json.Marshal(workerBody)
	req, _ = http.NewRequest("POST", "/api/worker/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid assign",
			body: map[string]interface{}{"repair_id": 1, "worker_id": 1},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid worker",
			body: map[string]interface{}{"repair_id": 1, "worker_id": 999},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/repair/assign", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestWorkerAcceptOrder(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	registerBody := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	repairBody := map[string]interface{}{
		"user_id":     loginResp.UserID,
		"username":    "testuser",
		"repair_type": "水电维修",
		"description": "水龙头漏水",
	}
	jsonBody, _ = json.Marshal(repairBody)
	req, _ = http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	workerBody := map[string]string{"work_no": "W001", "name": "张三", "phone": "13800138001", "password": "123456", "skill_type": "水电"}
	jsonBody, _ = json.Marshal(workerBody)
	req, _ = http.NewRequest("POST", "/api/worker/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assignBody := map[string]interface{}{"repair_id": 1, "worker_id": 1}
	jsonBody, _ = json.Marshal(assignBody)
	req, _ = http.NewRequest("POST", "/api/repair/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid accept",
			body: map[string]interface{}{"repair_id": 1, "worker_id": 1},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong worker",
			body: map[string]interface{}{"repair_id": 1, "worker_id": 2},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/repair/accept", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}

func TestWorkerRejectOrder(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()
	router := setupTestRouter()

	registerBody := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var loginResp LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &loginResp)

	repairBody := map[string]interface{}{
		"user_id":     loginResp.UserID,
		"username":    "testuser",
		"repair_type": "水电维修",
		"description": "水龙头漏水",
	}
	jsonBody, _ = json.Marshal(repairBody)
	req, _ = http.NewRequest("POST", "/api/repair/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	workerBody := map[string]string{"work_no": "W001", "name": "张三", "phone": "13800138001", "password": "123456", "skill_type": "水电"}
	jsonBody, _ = json.Marshal(workerBody)
	req, _ = http.NewRequest("POST", "/api/worker/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assignBody := map[string]interface{}{"repair_id": 1, "worker_id": 1}
	jsonBody, _ = json.Marshal(assignBody)
	req, _ = http.NewRequest("POST", "/api/repair/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid reject",
			body: map[string]interface{}{"repair_id": 1, "worker_id": 1},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/repair/reject", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.Code)
			}
		})
	}
}
