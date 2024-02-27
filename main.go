package main

import (
	//"context"
	//"context"
	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gbubemi22/moni/src/database"
	"github.com/gbubemi22/moni/src/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = ""
	}

	userService := user.NewUserService(&user.UserRepository{userCollection})


	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Wellcome to faster_moni"})
		
	})


	//Serve the Swagger UI at the /swagger URL
	// router.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userController := user.NewUserController(userService)
	user.AuthRoutes(router, userController)

	router.Run(":" + port)
}
