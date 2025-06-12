package main

import (
	"github.com/burhangltekin/byfood/routes"
	"github.com/burhangltekin/byfood/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitDB()

	r := gin.Default()
	r.Use(gin.Logger())

	routes.SetupRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
