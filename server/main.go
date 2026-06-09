package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := filepath.Join(".", "property_work_order.db")
	if err := InitDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer CloseDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(CORSMiddleware())

	frontendDir := filepath.Join("..", "web", "dist")

	if _, err := os.Stat(frontendDir); err == nil {
		r.StaticFS("/assets", gin.Dir(filepath.Join(frontendDir, "assets"), false))

		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			if filepath.Ext(path) != "" {
				filePath := filepath.Join(frontendDir, path)
				if _, err := os.Stat(filePath); err == nil {
					c.File(filePath)
					return
				}
			}

			if path == "/" || strings.HasPrefix(path, "/admin") || strings.HasPrefix(path, "/user") || strings.HasPrefix(path, "/worker") {
				c.File(filepath.Join(frontendDir, "index.html"))
				return
			}

			c.String(404, "Not Found")
		})
	}

	r.POST("/api/register", Register)
	r.POST("/api/login", Login)
	r.POST("/api/repair/create", CreateRepair)
	r.GET("/api/repair/user", GetUserRepairs)
	r.GET("/api/repair/all", GetAllRepairs)
	r.PUT("/api/repair/status", UpdateRepairStatus)
	r.POST("/api/repair/assign", AssignWorker)
	r.GET("/api/repair/worker", GetWorkerRepairs)
	r.POST("/api/repair/accept", WorkerAcceptOrder)
	r.POST("/api/repair/reject", WorkerRejectOrder)
	r.POST("/api/repair/result", SubmitRepairResult)
	r.GET("/api/repair/stats", GetWorkerStats)

	r.POST("/api/worker/login", WorkerLogin)
	r.GET("/api/worker/list", GetWorkers)
	r.POST("/api/worker/create", CreateWorker)
	r.PUT("/api/worker/update", UpdateWorker)
	r.PUT("/api/worker/status", ToggleWorkerStatus)

	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
