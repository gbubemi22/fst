package user

import (
	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up authentication routes
// AuthRoutes sets up user authentication routes
func AuthRoutes(incomingRoutes *gin.Engine, uc *UserController) {
	incomingRoutes.POST("/v1/auth/signup", uc.CreateUser)
	incomingRoutes.POST("/v1/auth/login", uc.Login)
	incomingRoutes.GET("/v1/users", uc.ListAllUsers)
	incomingRoutes.GET("/v1/users/:user_id", uc.ListOne)
}
