package main

import (
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/db"
	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/routes"
)

func main() {
	// Initialize the database connection
	db.InitDB()

	// Set up the router and start the server
	r := routes.SetupRouter()
	r.Run(":8080")
}
