package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/golang_gin_gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Product{})

	r := gin.Default()

	r.POST("/products", func(ctx *gin.Context) {
		var input Product

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
		}

		db.Create(&input)
		ctx.JSON(200, input)
	})

	r.GET("/products", func(ctx *gin.Context) {
		var products []Product
		db.Find(&products)

		ctx.JSON(200, products)
	})

	r.GET("/products/:id", func(ctx *gin.Context) {
		var product Product

		if err := db.First(&product, ctx.Param("id")).Error; err != nil {
			ctx.JSON(200, gin.H{"error": "Product Not Found"})
			return
		}

		ctx.JSON(200, product)
	})

	r.PUT("/products/:id", func(ctx *gin.Context) {
		var product Product

		if err := db.First(&product, ctx.Param("id")).Error; err != nil {
			ctx.JSON(200, gin.H{"error": "Product Not Found"})
			return
		}

		var input Product

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Model(&product).Updates(input)

		ctx.JSON(200, product)

	})

	r.DELETE("/products/:id", func(ctx *gin.Context) {
		var product Product
		if err := db.First(&product, ctx.Param("id")).Error; err != nil {
			ctx.JSON(200, gin.H{"error": "Product Not Found"})
			return
		}

		db.Delete(&product)
		ctx.JSON(200, gin.H{"message": "Product Deleted"})
	})

	r.Run(":8080")

}
