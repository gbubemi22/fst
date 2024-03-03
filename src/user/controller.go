package user

import (
	"context"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// Handler for creating a new user
func (uc *UserController) CreateUser(c *gin.Context) {
	// Create a context with timeout for the operation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Declare a variable to hold the new user data
	var newUser User

	// Bind JSON data from the request to newUser
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service method to create the user
	if err := uc.UserService.CreateUser(ctx, &newUser); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the created user in the response
	c.JSON(http.StatusCreated, newUser)
}

func (uc *UserController) Login(c *gin.Context) {
	var loginRequest map[string]interface{}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	email, emailOk := loginRequest["email"].(string)
	password, passwordOk := loginRequest["password"].(string)

	if !emailOk || !passwordOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Call the service to perform login
	user, token, err := uc.UserService.Login(context.Background(), email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Successful login
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": token})
}

// Handler for getting a user by ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	ctx := context.Background()

	user, err := uc.UserService.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler for getting a user by email
func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	ctx := context.Background()

	user, err := uc.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler for updating a user
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	if err := uc.UserService.UpdateUser(ctx, userID, &updatedUser); err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// Handler for deleting a user
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	ctx := context.Background()

	if err := uc.UserService.DeleteUser(ctx, userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (uc *UserController) ListAllUsers(c *gin.Context) {
	users, err := uc.UserService.ListAllUsers(c.Request.Context())
	if err != nil {
		if err.Error() == "no users found" {
			// Return a 404 status
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
			return
		}
		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the list of users
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (uc *UserController) ListOne(c *gin.Context) {
	userID := c.Param("id") // Assuming you extract the user ID from the request URL

	user, err := uc.UserService.ListOne(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			// Return a 404 status if the user is not found
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the user details
	c.JSON(http.StatusOK, gin.H{"user": user})
}
