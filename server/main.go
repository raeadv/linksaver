package main

import (
	"context"
	"fmt"
	"linksaver/server/database"
	"linksaver/server/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	goc := context.Background()

	fmt.Println("Loading env vars")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("establishing db connection")

	mode := os.Getenv("MODE")

	var db database.Connection

	// initiate gorm connection
	gormDb, err := db.InitGorm(goc)
	if err != nil {
		fmt.Println("Failed to initiate GORM")
		fmt.Printf("ERROR : %#v/n ", err)
		panic("gorm connection fail.")
	}

	database.Init(gormDb)

	fmt.Println("Database connection successfull")

	g := gin.Default()

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	if mode == "production" {
		g.Static("/assets", "./compiled/assets")
		g.StaticFile("/favicon.ico", "./compiled/favicon.ico")
		g.StaticFile("/", "./compiled/index.html")
		g.NoRoute(func(c *gin.Context) {
			c.File("./compiled/index.html")
		})
	}

	// register api routes
	api := g.Group("/api")

	routes.RegisterApiRoute(api)

	addr := os.Getenv("API_URL")
	if mode == "production" {
		addr = "0.0.0.0:8080"
	}

	g.Run(addr) // listens on 0.0.0.0:8080 by default
}
