package main

import (
	//"context"
	//"context"
	//"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/gbubemi22/moni/src/database"
	"github.com/gbubemi22/moni/src/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	// Set up MongoDB client
	mongoURI := os.Getenv("MONGO_URI")
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Set up user repository
	userRepository := user.NewUserRepository(mongoClient, "user")

	// Set up user service
	userService := user.NewUserService(userRepository)

	// Set up user controller
	userController := user.NewUserController(userService)

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Wellcome to faster_moni"})
	})

	user.AuthRoutes(router, userController)

	// userService := user.NewUserService()
	// userController := user.NewUserController(userService)
	// user.AuthRoutes(router, userController)

	router.Run(":" + port)
}
