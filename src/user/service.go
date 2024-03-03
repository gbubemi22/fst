package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	UserRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

type SignedDetails struct {
	Email      string
	// First_name string
	// Last_name  string
	// Phone      string
	// User_type  string
	jwt.RegisteredClaims
}

func (us *UserService) GenerateToken(email string) (string, error) {
	signedDetails := &SignedDetails{
		Email:      email,
		// First_name: firstName,
		// Last_name:  lastName,
		// Phone:      phone,
		// User_type:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, signedDetails).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	return token, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *User) error {
	// Check if the user with the given email or phone number already exists
	filter := bson.M{"$or": []bson.M{{"email": *user.Email}, {"phone": *user.Phone}}}
	cursor, err := us.UserRepository.Collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error checking for user existence: %v", err)
		return errors.New("error occurred while checking for the user")
	}
	defer cursor.Close(ctx)

	// Iterate over the results to check for email and phone existence
	for cursor.Next(ctx) {
		var existingUser User
		if err := cursor.Decode(&existingUser); err != nil {
			log.Printf("Error decoding user: %v", err)
			return errors.New("error occurred while decoding user")
		}
		if existingUser.Email != nil && *existingUser.Email == *user.Email {
			return errors.New("user with the provided email already exists")
		}

		if existingUser.Phone != nil && *existingUser.Phone == *user.Phone {
			return errors.New("user with the provided phone number already exists")
		}
	}

	// Hash the password
	password := HashPassword(*user.Password)
	user.Password = &password

	// Call the repository to create the user
	return us.UserRepository.CreateUser(ctx, user)
}
func (us *UserService) Login(ctx context.Context, email, password string) (*User, string, error) {
	// Get the user by email
	user, err := us.UserRepository.GetUserByEmail(ctx, email)
	fmt.Println("Check:::", user)
	if err != nil {
	    log.Printf("Error getting user by email: %v", err)
	    return nil, "", errors.New("invalid email or password")
	}
  
	if user == nil {
	    return nil, "", errors.New("user not found")
	}
  
	
  
	// Compare the provided password with the hashed password in the database
	if user.Password == nil || *user.Password == "" {
	    return nil, "", errors.New("user password is empty or nil")
	}
  
	check, msg := VerifyPassword(*user.Password, password)
	if !check {
	    return nil, "", errors.New(msg)
	}
  
	token, err := us.GenerateToken(*user.Email)
	if err != nil {
	    log.Printf("Error generating tokens: %v", err)
	    return nil, "", errors.New("failed to generate tokens")
	}
  
	// Passwords match, login successful
	return user, token, nil
  }
  
func (us *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	// Call the repository to get the user by ID
	return us.UserRepository.GetUserByID(ctx, userID)
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Call the repository to get the user by email
	return us.UserRepository.GetUserByEmail(ctx, email)
}

func (us *UserService) UpdateUser(ctx context.Context, userID string, updatedUser *User) error {
	// Perform any additional business logic or validation before calling the repository

	// Call the repository to update the user
	return us.UserRepository.UpdateUser(ctx, userID, updatedUser)
}

func (us *UserService) DeleteUser(ctx context.Context, userID string) error {
	// Perform any additional business logic or validation before calling the repository

	// Call the repository to delete the user
	return us.UserRepository.DeleteUser(ctx, userID)
}

func (us *UserService) ListAllUsers(ctx context.Context) ([]*User, error) {
	users, err := us.UserRepository.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		// No users found, return a 404 status
		return nil, errors.New("no users found")
	}

	return users, nil
}

func (us *UserService) ListOne(ctx context.Context, userId string) (*User, error) {
	user, err := us.UserRepository.ListOne(ctx, userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
