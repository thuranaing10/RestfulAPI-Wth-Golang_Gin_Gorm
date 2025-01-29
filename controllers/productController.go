package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/db"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/models"
)

// type Product struct {
// 	gorm.Model
// 	Name   string  `json:"name"`
// 	Price  float64 `json:"price"`
// 	UserID uint    `json:"user_id"`
// }

func CreateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{Name: input.Name, Price: input.Price, UserID: userID}
	if err := db.GetDB().Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts - Retrieves all products
func GetProducts(ctx *gin.Context) {
	var products []models.Product

	userID := ctx.MustGet("userID").(uint)

	if err := db.GetDB().Where("user_id = ?", userID).Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProduct - Retrieves a product by ID
func GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product

	userID := ctx.MustGet("userID").(uint)

	if err := db.GetDB().Where("user_id = ?", userID).First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product": product})
}

// UpdateProduct - Updates a product by ID
func UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product

	userID := ctx.MustGet("userID").(uint)

	if err := db.GetDB().Where("user_id = ?", userID).First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.GetDB().Save(&product)
	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}

// DeleteProduct - Deletes a product by ID
func DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product
	userID := ctx.MustGet("userID").(uint)

	if err := db.GetDB().Where("user_id = ?", userID).First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	db.GetDB().Delete(&product)
	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
