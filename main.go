package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Secret key for signing the JWT (store securely in environment variables)
var jwtSecret = []byte("your-secret-key")

// JWT claims structure
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Products []Product `json:"products" gorm:"foreignKey:UserID"`
}

type Product struct {
	gorm.Model
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	UserID uint    `json:"user_id"`
}

var db *gorm.DB

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/golang_gin_gorm?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&User{}, &Product{})

	r := gin.Default()

	r.POST("/register", register)
	r.POST("/login", login)

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/products", createProduct)
		auth.GET("/products", getProducts)
		auth.GET("/products/:id", getProduct)
		auth.PUT("/products/:id", updateProduct)
		auth.DELETE("/products/:id", deleteProduct)
	}

	// Start the server
	r.Run(":8080")

}

func register(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	if err := db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return // Ensure you return after handling the error
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

// Login function with JWT generation
func login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user in the database
	var user User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return JWT token
	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// AuthMiddleware validates JWT token and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Check if token follows "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user info in context for later use
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// Create a product
func createProduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found"})
		return
	}

	// Ensure correct type assertion
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UserID type"})
		return
	}

	var input struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := Product{
		Name:   input.Name,
		Price:  input.Price,
		UserID: userIDUint,
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Get all products for the logged-in user
func getProducts(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var products []Product
	if err := db.Where("user_id = ?", userID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Get a single product by ID (only if it belongs to the logged-in user)
func getProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var product Product
	if err := db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update a product (only if it belongs to the logged-in user)
func updateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the product
	var product Product
	if err := db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Update the product
	db.Model(&product).Updates(input)
	c.JSON(http.StatusOK, product)
}

// Delete a product (only if it belongs to the logged-in user)
func deleteProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Find the product
	var product Product
	if err := db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Delete the product
	db.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
