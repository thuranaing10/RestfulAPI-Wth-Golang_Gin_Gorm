package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/controllers"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/products")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/", controllers.CreateProduct)
		auth.GET("/", controllers.GetProducts)
		auth.GET("/:id", controllers.GetProduct)
		auth.PUT("/:id", controllers.UpdateProduct)
		auth.DELETE("/:id", controllers.DeleteProduct)
	}

	return r
}
