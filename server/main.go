package main

import (
	"log"
	"os"
	"path/filepath"

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
		r.GET("/admin/*any", serveIndex(frontendDir))
		r.GET("/user/*any", serveIndex(frontendDir))
		r.Static("/assets", filepath.Join(frontendDir, "assets"))
		r.GET("/", func(c *gin.Context) {
			c.File(filepath.Join(frontendDir, "index.html"))
		})
	}

	r.POST("/api/register", Register)
	r.POST("/api/login", Login)
	r.POST("/api/repair/create", CreateRepair)
	r.GET("/api/repair/user", GetUserRepairs)
	r.GET("/api/repair/all", GetAllRepairs)
	r.PUT("/api/repair/status", UpdateRepairStatus)

	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func serveIndex(frontendDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			c.String(404, "Frontend not found")
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
