package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/db"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/models"
)

func CreatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	post := models.Post{Title: input.Title, Description: input.Description, UserID: userID}

	if err := db.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post!",
		})

		return
	}

	c.JSON(http.StatusOK, post)

}

func GetPosts(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var posts []models.Post

	if err := db.GetDB().Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch posts",
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var post models.Post

	if err := db.GetDB().Where("user_id = ?", userID).Find(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})

}

func UpdatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var post models.Post

	if err := db.GetDB().Where("user_id = ?", userID).First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found!",
		})

		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	db.GetDB().Save(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"post":    post,
	})

}

func DeletePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")
	var post models.Post

	if err := db.GetDB().Where("user_id = ?", userID).First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found!",
		})

		return
	}

	db.GetDB().Delete(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})

}
