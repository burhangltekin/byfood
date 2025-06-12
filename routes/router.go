package routes

import (
	"github.com/burhangltekin/byfood/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/books", controllers.GetBooks)
		api.GET("/books/:id", controllers.GetBook)
		api.POST("/books", controllers.CreateBook)
		api.PUT("/books/:id", controllers.UpdateBook)
		api.DELETE("/books/:id", controllers.DeleteBook)
	}
}
