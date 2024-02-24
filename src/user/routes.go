package user

import (
	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up authentication routes
// AuthRoutes sets up user authentication routes
func AuthRoutes(incomingRoutes *gin.Engine, uc *UserController) {
	v1 := incomingRoutes.Group("/v1") // Add versioning
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/signup", uc.CreateUser)
			authGroup.POST("/login", uc.Login)
		}

		userGroup := v1.Group("/users")

		{
			userGroup.GET("/", uc.ListAllUsers)
			userGroup.GET("/:user_id", uc.ListOne) // Adjust the function name
		}
	}

	// You can add more routes or versions as needed
}
