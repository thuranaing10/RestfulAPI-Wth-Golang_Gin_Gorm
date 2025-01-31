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

	// auth := r.Group("/products")
	r.Use(middlewares.AuthMiddleware())
	{
		product := r.Group("/products")
		product.POST("/", controllers.CreateProduct)
		product.GET("/", controllers.GetProducts)
		product.GET("/:id", controllers.GetProduct)
		product.PUT("/:id", controllers.UpdateProduct)
		product.DELETE("/:id", controllers.DeleteProduct)

		post := r.Group("/posts")
		post.POST("/", controllers.CreatePost)
		post.GET("/", controllers.GetPosts)
		post.GET("/:id", controllers.GetPost)
		post.PUT("/:id", controllers.UpdatePost)
		post.DELETE("/:id", controllers.DeletePost)

	}

	return r
}
