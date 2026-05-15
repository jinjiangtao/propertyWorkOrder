package router

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
		api.POST("/admin/login", controller.AdminLogin)

		api.POST("/workorder", controller.CreateWorkOrder)
		api.GET("/workorders", controller.GetWorkOrders)
		api.PUT("/workorder/:id", controller.UpdateWorkOrderStatus)
	}

	return r
}