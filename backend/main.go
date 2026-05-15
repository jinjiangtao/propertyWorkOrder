package main

import (
	"backend/model"
	"backend/router"
	"log"
)

func main() {
	err := model.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := router.SetupRouter()

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}